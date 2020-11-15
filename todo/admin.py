from todo.models import Project
from django.contrib import admin

# Register your models here.
from . import models

# admin.site.register(models.Project)
admin.site.register(models.Log)


@admin.register(Project)
class ProjectAdmin(admin.ModelAdmin):
    date_hierarchy = 'createdAt'
    pass
