var onboarding = {
    currentStep: parseInt(document.body.dataset.onboardingStep || '7'),
    active: document.body.dataset.onboardingActive === 'True',
    overlay: null,
    coachBody: null,
    coachContent: null,
    nextBtn: null,
    backBtn: null,
    skipBtn: null,
    stepBadge: null,
    stepTitle: null,
    pollInterval: null,
    allowed: false,

    steps: [
        { id: 1, title: 'Bienvenido', target: null, skip: false,
          html: '<div class="text-center"><div class="text-3xl mb-3">💰</div><h3 class="text-base font-bold text-chele-text mb-2">Bienvenido a Chele</h3><p class="text-sm text-chele-text-secondary">El método de los sobres: cada peso en tus cuentas debe tener un trabajo. Vas a asignar tu dinero a categorías.</p></div>' },
        { id: 2, title: 'Tus cuentas', target: 'a[href*="/cuentas/crear/"]', skip: true,
          html: '<p><strong>Agregá tus cuentas bancarias.</strong> El saldo inicial aparece en "Por asignar" — lo vas a distribuir después.</p>' },
        { id: 3, title: 'Asignar dinero', target: '[class*="Por asignar"]', skip: false,
          html: '<p><strong>Este es tu dinero sin trabajo.</strong> Asignalo a categorías usando los inputs de cada una. El objetivo es que "Por asignar" llegue a <strong>$0</strong>.</p><div id="rtaProgress" class="mt-3 bg-chele-bg rounded-lg p-3 text-center"><span class="text-lg font-bold text-chele-warning">Cargando...</span></div>' },
        { id: 4, title: 'Metas de ahorro', target: '.category-item, tr.cursor-pointer', skip: true,
          html: '<p><strong>¿Querés ahorrar para algo?</strong> Hacé click en una categoría y agregá una meta. Te dice cuánto asignar cada mes.</p>' },
        { id: 5, title: 'Ingresos fijos', target: 'a[href*="/programaciones/"]', skip: true,
          html: '<p><strong>¿Tenés gastos o ingresos que se repiten?</strong> Programalos y Chele los aplica automáticamente.</p>' },
        { id: 6, title: 'Tus deudas', target: 'a[href*="/tc/crear/"], a[href*="/prestamos/crear/"]', skip: true,
          html: '<p><strong>Si tenés tarjetas o préstamos, agregalos.</strong> Chele maneja la deuda automáticamente.</p>' },
        { id: 7, title: '¡Todo listo!', target: null, skip: false,
          html: '<div class="text-center"><div class="text-3xl mb-3">🎉</div><h3 class="text-base font-bold text-chele-text mb-2">¡Todo listo!</h3><p class="text-sm text-chele-text-secondary mb-3">Tu dinero ya tiene trabajo. Recordá registrar tus gastos del día a día y revisar tus categorías cada mes.</p></div>' }
    ],

    init: function() {
        if (!this.active) return;
        this.overlay = document.getElementById('onboardingOverlay');
        this.coachBody = document.getElementById('onboardingCoachBody');
        this.coachContent = document.getElementById('onboardingCoachContent');
        this.nextBtn = document.getElementById('onboardingNextBtn');
        this.backBtn = document.getElementById('onboardingBackBtn');
        this.skipBtn = document.getElementById('onboardingSkipBtn');
        this.stepBadge = document.getElementById('onboardingStepBadge');
        this.stepTitle = document.getElementById('onboardingStepTitle');

        var self = this;
        this.nextBtn.onclick = function() { self.advance(); };
        this.backBtn.onclick = function() { self.goTo(self.currentStep - 1); };
        this.skipBtn.onclick = function() { self.skip(); };
        this.overlay.style.display = 'flex';
        this.goTo(this.currentStep);

        if (this.currentStep < 7) {
            this.pollInterval = setInterval(function() { self.poll(); }, 3000);
        }
    },

    goTo: function(step) {
        if (step < 1 || step > 7) return;
        this.currentStep = step;
        var s = this.steps[step - 1];
        this.stepBadge.textContent = step + '/7';
        this.stepTitle.textContent = s.title;
        this.coachContent.innerHTML = s.html;

        this.backBtn.style.display = step > 1 ? '' : 'none';
        this.skipBtn.style.display = (s.skip && step < 7) ? '' : 'none';
        this.nextBtn.textContent = step === 7 ? 'Ir al presupuesto' : 'Siguiente';
        this.nextBtn.disabled = !this.allowed;
        this.nextBtn.className = 'text-xs bg-chele-primary text-white rounded-lg px-4 py-2 font-medium hover:bg-chele-primary-dark disabled:opacity-40 disabled:cursor-not-allowed';

        this.updateProgress(step);
        this.highlightTarget(s.target);
        this.poll();
    },

    advance: function() {
        if (this.nextBtn.disabled) return;
        if (this.currentStep >= 7) {
            this.complete();
            return;
        }
        var self = this;
        fetch('/onboarding/avanzar/', { method: 'POST', headers: { 'X-CSRFToken': this.csrf() } })
        .then(function(r) { return r.json(); })
        .then(function(data) {
            self.currentStep = data.step;
            self.goTo(data.step);
            self.allowed = false;
            self.nextBtn.disabled = true;
        });
    },

    complete: function() {
        var self = this;
        fetch('/onboarding/avanzar/', { method: 'POST', headers: { 'X-CSRFToken': this.csrf() } })
        .then(function() { self.overlay.style.display = 'none'; location.href = '/presupuestos/'; });
    },

    skip: function() {
        this.allowed = true;
        this.advance();
    },

    poll: function() {
        var self = this;
        fetch('/onboarding/state/')
        .then(function(r) { return r.json(); })
        .then(function(data) {
            if (data.step_completed) {
                self.allowed = true;
                self.nextBtn.disabled = false;
                self.nextBtn.className = 'text-xs bg-chele-success text-white rounded-lg px-4 py-2 font-medium hover:bg-chele-success/80 animate-pulse';
            }
            if (self.currentStep === 3 && data.step_completed) {
                document.getElementById('rtaProgress').innerHTML = '<span class="text-lg font-bold text-chele-success">✅ $0 — Todo tiene trabajo</span>';
            } else if (self.currentStep === 3) {
                document.getElementById('rtaProgress').innerHTML = '<span class="text-lg font-bold text-chele-warning">$' + data.ready_to_assign.toFixed(2) + ' sin asignar</span>';
            }
        });
    },

    highlightTarget: function(selector) {
        document.querySelectorAll('.onboarding-highlight').forEach(function(el) {
            el.classList.remove('onboarding-highlight');
        });
        if (!selector) return;
        var el = document.querySelector(selector);
        if (el) {
            el.classList.add('onboarding-highlight');
            el.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    },

    updateProgress: function(step) {
        document.querySelectorAll('#onboardingOverlay [data-step]').forEach(function(dot) {
            var s = parseInt(dot.dataset.step);
            dot.className = 'h-1.5 w-6 rounded-full transition-colors duration-300 ' + (s <= step ? 'bg-chele-success' : 'bg-chele-border');
        });
    },

    csrf: function() {
        return document.querySelector('[name=csrfmiddlewaretoken]')?.value || '';
    }
};

document.addEventListener('DOMContentLoaded', function() { onboarding.init(); });
