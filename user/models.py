from todo.views import index
from django.db import models
from django.forms import PasswordInput
# Create your models here.


class User(models.Model):
    nickName = models.CharField(max_length=12, unique=True, null=False)
    password = models.CharField(
        max_length=128, null=False)
    email = models.EmailField(unique=True)
    last_login = models.DateTimeField(auto_now=True)
    create_at = models.DateTimeField(auto_now_add=True)

    def __str__(self) -> str:
        return self.nickName + '[-create by]' + self.create_at.strftime("%Y-%m-%d")
