from django.urls import path

from . import views
urlpatterns = [
    path("", views.index, name='index'),
    path("detail/<int:project_id>", views.detail, name='detail'),
    path("detail/add", views.add, name="add"),
    path("detail/<int:project_id>/del", views.delete, name='delete'),
    path("detail/<int:project_id>/edit", views.update, name="update"),

    path("detail/<int:project_id>/logs/add", views.addLog, name='addLog'),
    path("detail/<int:project_id>/logs/<int:log_id>/edit",
         views.EditLog, name="editLog"),

    path("hello", views.hello, name='hello')
]
