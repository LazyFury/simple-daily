from django.db.models import fields
from django.forms import ModelForm
from .models import Log, Project


class ProjectForm(ModelForm):
    class Meta:
        model = Project
        fields = ['name', 'progress', 'describe']


class AddProjectForm(ModelForm):
    class Meta:
        model = Project
        fields = ['name', 'progress', 'describe',
                  'start', 'expectedEnd', 'actualDeliveryDate']


class ProjectLogForm(ModelForm):
    class Meta:
        model = Log
        fields = ["content", "plusProgress"]
