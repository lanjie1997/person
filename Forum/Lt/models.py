from __future__ import unicode_literals
from django.db import models
from django.contrib.auth.models import User

# Create your models here.



class Article(models.Model):
    '''
    帖子表
    '''
    title = models.CharField(u'文章标题',max_length=255,unique=True)
    content = models.TextField(u"内容")
    author = models.ForeignKey("UserProfile",verbose_name="作者")
    publish_date = models.DateTimeField(auto_now=True,verbose_name="发布日期")
    # priority = models.IntegerField(default=1000,verbose_name="优先级")

    def __unicode__(self):
        return "<%s,author:%s>" % (self.title,self.author)

class Comment(models.Model):
    '''
    评论表
    '''
    article = models.ForeignKey("Article")
    user = models.ForeignKey("UserProfile")
    comment = models.TextField(max_length=1000)
    date = models.DateTimeField(auto_now=True)
    # parent_comment = models.ForeignKey("self",related_name='p_comment',blank=True,null=True)

    def __unicode__(self):
        return "<user:%s>" %(self.user)
# class ThumbUp(models.Model):
#     '''
#     点赞
#     '''
#     article = models.ForeignKey('Article')
#     user = models.ForeignKey('UserProfile')
#     date = models.DateTimeField(auto_now=True)


class UserProfile(models.Model):
    '''
    用户表
    '''
    username = models.CharField(max_length=30)
    password = models.CharField(max_length=50)
    email = models.EmailField(max_length=50)
    sex = models.CharField(max_length=10)
    # last_login = models.DateTimeField(blank=True)

    def __unicode__(self):
        return self.name
class UserGroup(models.Model):
    '''
    用户组表
    '''
    name = models.CharField(max_length=64,unique=True)
    def __unicode__(self):
        return self.name
