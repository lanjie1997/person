from django.contrib import admin
from .models import *
# Register your models here.

class ArticleAdmin(admin.ModelAdmin):
    list_display = ('id','title','author','publish_date')

admin.site.register(Article,ArticleAdmin)
admin.site.register(Comment)
admin.site.register(UserProfile)
admin.site.register(UserGroup)