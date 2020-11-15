from datetime import date
from user.utils import needUser
from django.core.paginator import Paginator
from .form import AddProjectForm, ProjectForm, ProjectLogForm
from django.utils import timezone
from django.http.request import HttpRequest
from django.http.response import Http404, HttpResponse
from .models import Log, Project
import datetime
from django.contrib.auth.decorators import login_required
from django.shortcuts import get_list_or_404, get_object_or_404, redirect, render
# Create your views here.


def hello(request):
    return HttpResponse("hello world! django wsgi easy restart")


@needUser
@login_required
def index(request: HttpRequest):
    print(request.uid)
    if request.method == "GET":
        page = int(request.GET.get("page", default=1))
        pageSize = int(request.GET.get("size", default=10))
        if(page < 1):
            page = 1
        if(pageSize < 1):
            pageSize = 1
        projects = Project.objects.all().order_by(
            "-updatedAt")
        paginator = Paginator(projects, pageSize)
        try:
            projectList = paginator.page(page)
        except:
            projectList = []
        context = {
            "projects": projectList,
            "size": pageSize
        }
        return render(request, "index.html", context)


def add(request: HttpRequest):
    if(request.method == "GET"):
        return render(request, "add.html")
    if(request.method == "POST"):
        project = Project()
        project.createdAt = timezone.now()
        project.updatedAt = timezone.now()
        p = AddProjectForm(data=request.POST, instance=project)
        if(p.is_valid()):
            p.save()
            return redirect("index")
        return HttpResponse("error")


def update(request, project_id):
    p = get_object_or_404(Project, pk=project_id)
    if(request.method == "GET"):
        form = ProjectForm(instance=p)
        context = {
            "project": p,
            "form": form,
            "now": timezone.now()
        }
        return render(request, 'edit.html', context)
    if(request.method == "POST"):
        p.updatedAt = timezone.now()
        form = ProjectForm(instance=p, data=request.POST)
        form.save()
        if(form.is_valid()):
            return redirect("index")
        else:
            return HttpResponse("error")


def delete(request, project_id):
    project = get_object_or_404(Project, pk=project_id)
    try:
        project.delete()
        return redirect("index")
    except:
        return Http404("删除失败")


def detail(request: HttpRequest, project_id):
    p = get_object_or_404(Project, pk=project_id)
    now = datetime.datetime.now()
    start = datetime.date(1970, 1, 1)
    type = request.GET.get("type")
    if(type != ""):
        delta = {
            "day": datetime.timedelta(hours=now.hour),
            "yesterday": datetime.timedelta(hours=now.hour+24),
            "week": datetime.timedelta(days=now.weekday()),
            "month": datetime.timedelta(days=now.day)
        }.get(type, None)

        print(delta)
        if(delta != None):
            start = now - delta
    logs = Log.objects.filter(project=Project(
        id=project_id), createdAt__gte=start).order_by("-createdAt")
    plus_progress = 0
    for l in logs.all():
        plus_progress += l.plusProgress
    context = {
        "project": p,
        "logs": logs,
        "jobs": ",".join(l.content for l in logs),
        "isWeek": type == 'week',
        "plus_progress": plus_progress
    }
    return render(request, "detail.html", context)


def addLog(request: HttpRequest, project_id: int):
    project = get_object_or_404(Project, pk=project_id)
    if(request.method == "GET"):
        now = timezone.now().strftime("%Y-%m-%dT%H:%I")
        return render(request, "logs/add.html", {
            "project": project,
            "now": now
        })
    if(request.method == "POST"):
        log = Log(project=project, createdAt=timezone.now(),
                  updatedAt=timezone.now())
        form = ProjectLogForm(instance=log, data=request.POST)
        _log: Log = form.save()
        if(form.is_valid()):
            project.progress += _log.plusProgress
            try:
                project.save()
                return redirect("/detail/%d" % project.id)
            except:
                pass
        return HttpResponse("提交失败")


def EditLog(request: HttpRequest, project_id, log_id):
    project = get_object_or_404(Project, pk=project_id)
    log: Log = get_object_or_404(Log, pk=log_id)
    if(request.method == "GET"):
        return render(request, "logs/edit.html", {
            "project": project,
            "log": log
        })
    if(request.method == "POST"):
        project.progress -= log.plusProgress

        form = ProjectLogForm(data=request.POST, instance=log)
        _log = form.save()
        if(form.is_valid()):
            project.progress += _log.plusProgress
            try:
                project.save()
                return redirect("/detail/%d" % project.id)
            except:
                pass
        return HttpResponse("更新失败")
