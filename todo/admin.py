from django.contrib import admin

# Register your models here.
from . import models

admin.site.register(models.Project)
admin.site.register(models.Log)
