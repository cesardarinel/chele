from django.db import models
from django.conf import settings


class OnboardingProfile(models.Model):
    user = models.OneToOneField(settings.AUTH_USER_MODEL, on_delete=models.CASCADE, related_name='onboarding')
    step = models.IntegerField(default=0, verbose_name='paso de onboarding')

    class Meta:
        verbose_name = 'onboarding'
        verbose_name_plural = 'onboardings'

    def __str__(self):
        return f'{self.user.username}: paso {self.step}'
