from django.contrib import admin
from .models import OnboardingProfile


@admin.register(OnboardingProfile)
class OnboardingProfileAdmin(admin.ModelAdmin):
    list_display = ('user', 'step')
    list_filter = ('step',)
    search_fields = ('user__username', 'user__email')
