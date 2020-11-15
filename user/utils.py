from functools import wraps
from django.http import HttpResponse
from django.http.request import HttpRequest


def needUser(fn) -> HttpResponse:
    def decorator(request: HttpRequest, *args, **kwargs):
        print(request.session.values)
        request.uid = 1
        return fn(request, *args, **kwargs)
    return decorator
