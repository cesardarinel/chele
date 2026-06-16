from django import template
from django.contrib.humanize.templatetags.humanize import intcomma

register = template.Library()


@register.filter
def currency(value):
    try:
        value = float(value)
    except (TypeError, ValueError):
        return value
    formatted = f'{value:,.2f}'
    return formatted
