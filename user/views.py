from user.forms import UserForm
from django.shortcuts import redirect, render
from django.http import HttpRequest, HttpResponse
# Create your views here.


def login(request: HttpRequest):
    return HttpResponse("login")


def reg(request: HttpRequest):
    if(request.method == "GET"):
        return render(request, "user/reg.html", {})
    if(request.method == "POST"):
        form = UserForm(data=request.POST)
        form.save()
        if(form.is_valid()):

            return redirect("index")
        return HttpResponse("reg")
        # return redirect("index")


def forgot(request: HttpRequest):
    return HttpResponse("forgot")
