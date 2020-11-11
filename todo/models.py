from django.db import models
from django.db.models.query import EmptyQuerySet
from django.utils import timezone
import datetime
# Create your models here.


class Project(models.Model):
    name = models.CharField(max_length=200)
    describe = models.CharField(max_length=500)
    start = models.DateField("start_time")
    expectedEnd = models.DateField("expected_end_time")
    actualDeliveryDate = models.DateField("actual_delivery_date")
    createdAt = models.DateTimeField("created_at")
    updatedAt = models.DateTimeField("updated_at")
    progress = models.IntegerField(default=0)

    def getStartTime(self):
        return self.start.strftime("%Y-%m-%d")

    def getExpectedEnd(self):
        return self.expectedEnd.strftime("%Y-%m-%d")

    def getActualDeliveryDate(self):
        return self.actualDeliveryDate.strftime("%Y-%m-%d")

    def getUpdatedAt(self):
        return self.updatedAt.strftime("%Y-%m-%dT%H:%I")

    def getCreatedAt(self):
        return self.createdAt.strftime("%Y-%m-%dT%H:%I")

    def __str__(self):
        return self.name


class Log(models.Model):
    # TODO:ForeignKey on_delete
    project = models.ForeignKey(Project, on_delete=models.CASCADE)
    content = models.TextField(max_length=500)
    plusProgress = models.IntegerField("plus_progress", default=0)
    createdAt = models.DateTimeField("created_at")
    updatedAt = models.DateTimeField("updated_at")

    def getCreatedAt(self):
        return self.createdAt.strftime("%Y-%m-%d %H:%I:%S")

    def getEditCreatedAt(self):
        return self.createdAt.strftime("%Y-%m-%dT%H:%I")

    def isAddToday(self):
        now = timezone.now()
        delta = datetime.timedelta(hours=datetime.datetime.now().hour)
        return now - self.createdAt <= delta

    def __str__(self):
        return self.content
