from django.urls import path
from .views import forgot, login, reg
urlpatterns = [
    path("login", login, name='user.login'),
    path("reg", reg, name='user.reg'),
    path("forgot", forgot, name='user.forgot')
]
