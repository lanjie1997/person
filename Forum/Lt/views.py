from django.shortcuts import render
from django.shortcuts import HttpResponseRedirect,render,HttpResponse,redirect
from django.contrib import auth
from django import forms
from django.contrib.auth.models import User
from .models import *
import datetime
# Create your views here.
def index(request):
    Forum=Article.objects.all()
    return  render(request,'index.html',{'Forum':Forum})

class RegistForm(forms.Form):
    gender = (
        ('male', "男"),
        ('female', "女"),
    )
    username = forms.CharField(label='用户名')
    password = forms.CharField(label='密   码', widget=forms.PasswordInput())
    sex = forms.ChoiceField(label='性别', choices=gender)
    email=forms.EmailField(label='邮 箱')
    # last_login = forms.DateTimeField()
class UserForm(forms.Form):
    username = forms.CharField(label='用户名')
    password = forms.CharField(label='密   码', widget=forms.PasswordInput())

class ChangeForm(forms.Form):
    username = forms.CharField(label='用户名')
    old_password = forms.CharField(label='原密码', widget=forms.PasswordInput())
    new_password = forms.CharField(label='新密码', widget=forms.PasswordInput())



def regist(request):
    if request.method == 'POST':
        uf = RegistForm(request.POST)
        if uf.is_valid():
            username = uf.cleaned_data['username']
            password = uf.cleaned_data['password']
            email=uf.cleaned_data['email']
            sex=uf.cleaned_data['sex']
            ##判断用户原密码是否匹配
            user = UserProfile.objects.filter(username=username)
            if user:
                return HttpResponse("用户名已存在")
            elif len(user) == 0:
                user = UserProfile()
                user.username = username
                user.password = password
                user.email=email
                user.sex=sex
                user.save()

            return HttpResponseRedirect('/login/')
    else:
        uf = RegistForm()

    return render(request,'regist.html', {'uf': uf})


def login(request):
    if request.session.get('is_login',None):
        return redirect('/index/')
    if request.method == 'POST':
        ##获取表单信息
        uf = UserForm(request.POST)
        if uf.is_valid():
            username = uf.cleaned_data['username']
            password = uf.cleaned_data['password']
            user = UserProfile.objects.get(username=username)
            if user:
                passwd = UserProfile.objects.filter(username=username, password=password)
                if passwd:
                    request.session['is_login'] = True
                    request.session['user_id'] = user.id
                    request.session['user_name'] = user.username
                    return redirect('/index/')
                else:
                    return HttpResponse('请检查密码是否正确!')
            elif len(user) == 0:
                return HttpResponse ('请检查用户名是否正确!')
    else:
        uf = UserForm()
    return render(request,'login.html', {'uf': uf})


def change_pass(request):
    if request.method == 'POST':
        uf = ChangeForm(request.POST)
        if uf.is_valid():
            username = uf.cleaned_data['username']
            old_password = uf.cleaned_data['old_password']
            new_password = uf.cleaned_data['new_password']

            ##判断用户原密码是否匹配
            user = UserProfile.objects.filter(username=username)
            if user:
                passwd = UserProfile.objects.filter(username=username, password=old_password)
                if passwd:
                    UserProfile.objects.filter(username=username, password=old_password).update(
                        password=new_password)  ##如果用户名、原密码匹配则更新密码
                    info = '密码修改成功!'
                else:
                    info = '请检查原密码是否输入正确!'
            elif len(user) == 0:
                info = '请检查用户名是否正确!'

        return HttpResponse(info)
    else:
        uf = ChangeForm()
    return render(request,'change.html', {'uf': uf})

def logout(request):
    if not request.session.get('is_login', None):
        # 如果本来就未登录，也就没有登出一说
        return redirect("/index/")
    request.session.flush()
    return redirect("/index/")


def article(request):
    return render(request,'article.html')
def article_form(request):
    title=request.POST.get('title')
    author=UserProfile.objects.get(username=request.session['user_name'])
    content=request.POST.get('content')
    # category=request.POST.get('category')
    # hidden=request.POST.get('hidden')
    # publish_date=datetime.now()
    Article.objects.create(title=title,author=author,content=content)
    return HttpResponse("<h1>发帖成功, 点击返回 <br/><a href='/index/'>主页</a><h1>")

def tz(request,Forum_id):
    tz=Article.objects.get(id=Forum_id)
    return render(request,'tz.html',{'tz':tz})
