var cheleGuide = {
    active: false,
    condition: null,
    completed: false,
    overlay: null,
    coachBody: null,
    content: null,
    nextBtn: null,
    dismissBtn: null,
    stepBadge: null,
    pollInterval: null,

    conditions: {
        'welcome': {
            title: 'Bienvenido',
            target: null,
            icon: '💰',
            text: 'Bienvenido a Chele. Cada peso debe tener un trabajo. Vamos a organizar tu dinero paso a paso.',
            btnText: 'Empezar',
            btnAction: 'advance'
        },
        'accounts': {
            title: 'Tus cuentas',
            target: 'a[href*="/cuentas/crear/"]',
            icon: '🏦',
            text: 'Agregá tus cuentas bancarias. El saldo inicial lo vas a distribuir después.',
            btnText: 'Ya agregué',
            btnAction: 'advance'
        },
        'assign': {
            title: 'Dinero sin asignar',
            target: '[class*="Por asignar"], [class*="ready"]',
            icon: '📊',
            text: 'Tenés dinero sin trabajo. Asignalo a categorías hasta que "Por asignar" llegue a $0.',
            btnText: 'Cerrar',
            btnAction: 'dismiss'
        },
        'goals': {
            title: 'Metas de ahorro',
            target: '.category-item, tr.cursor-pointer',
            icon: '🎯',
            text: '¿Querés ahorrar para algo? Hacé click en una categoría y poné una meta.',
            btnText: 'Cerrar',
            btnAction: 'dismiss'
        },
        'schedules': {
            title: 'Gastos recurrentes',
            target: 'a[href*="/programaciones/"]',
            icon: '📅',
            text: '¿Tenés gastos o ingresos que se repiten? Programalos y se aplican solos.',
            btnText: 'Cerrar',
            btnAction: 'dismiss'
        }
    },

    init: function() {
        this.overlay = document.getElementById('cheleGuideOverlay');
        if (!this.overlay) return;
        this.coachBody = document.getElementById('cheleGuideBody');
        this.content = document.getElementById('cheleGuideContent');
        this.nextBtn = document.getElementById('cheleGuideNextBtn');
        this.dismissBtn = document.getElementById('cheleGuideDismissBtn');
        this.stepBadge = document.getElementById('cheleGuideBadge');

        var self = this;
        this.nextBtn.onclick = function() { self.doAction(); };
        this.dismissBtn.onclick = function() { self.hide(); };

        this.poll();
        this.pollInterval = setInterval(function() { self.poll(); }, 4000);
    },

    poll: function() {
        var self = this;
        fetch('/onboarding/state/')
        .then(function(r) { return r.json(); })
        .then(function(data) {
            if (data.active && data.condition) {
                self.show(data.condition, data.condition_completed, data.ready_to_assign);
            } else {
                self.hide();
            }
        });
    },

    show: function(conditionId, completed, rta) {
        this.condition = conditionId;
        this.completed = completed;
        var c = this.conditions[conditionId];
        if (!c) return;

        this.stepBadge.textContent = c.icon;
        var html = '<h3 class="text-base font-bold text-chele-text mb-2">' + c.title + '</h3>';
        html += '<p class="text-sm text-chele-text-secondary">' + c.text + '</p>';

        if (conditionId === 'assign') {
            html += '<div class="mt-3 bg-chele-bg rounded-lg p-3 text-center"><span class="text-lg font-bold ' + (rta > 0 ? 'text-chele-warning' : 'text-chele-success') + '">$' + rta.toFixed(2) + (rta > 0 ? ' sin asignar' : ' — Todo tiene trabajo') + '</span></div>';
        }

        this.content.innerHTML = html;
        this.nextBtn.textContent = completed ? '✅ Siguiente' : c.btnText;
        this.nextBtn.disabled = false;

        if (completed) {
            this.nextBtn.className = 'text-xs bg-chele-success text-white rounded-lg px-4 py-2 font-medium hover:bg-chele-success/80';
        } else if (conditionId === 'assign') {
            this.nextBtn.className = 'text-xs bg-chele-bg-tertiary text-chele-text rounded-lg px-4 py-2 font-medium opacity-60 cursor-not-allowed';
            this.nextBtn.disabled = true;
        } else {
            this.nextBtn.className = 'text-xs bg-chele-primary text-white rounded-lg px-4 py-2 font-medium hover:bg-chele-primary-dark';
        }

        this.highlightTarget(c.target);
        this.overlay.style.display = 'flex';
        this.active = true;
    },

    hide: function() {
        this.overlay.style.display = 'none';
        this.clearHighlight();
        this.active = false;
    },

    doAction: function() {
        if (this.nextBtn.disabled) return;
        var c = this.conditions[this.condition];
        if (c.btnAction === 'advance') {
            var self = this;
            fetch('/onboarding/avanzar/', { method: 'POST', headers: { 'X-CSRFToken': this.csrf() } })
            .then(function() { self.hide(); });
        } else {
            this.hide();
        }
    },

    highlightTarget: function(selector) {
        this.clearHighlight();
        if (!selector) return;
        var el = document.querySelector(selector);
        if (el) {
            el.classList.add('chele-guide-highlight');
            el.scrollIntoView({ behavior: 'smooth', block: 'center' });
        }
    },

    clearHighlight: function() {
        document.querySelectorAll('.chele-guide-highlight').forEach(function(el) {
            el.classList.remove('chele-guide-highlight');
        });
    },

    csrf: function() {
        return document.querySelector('[name=csrfmiddlewaretoken]')?.value || '';
    }
};

document.addEventListener('DOMContentLoaded', function() { cheleGuide.init(); });
