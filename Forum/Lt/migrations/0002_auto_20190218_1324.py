# -*- coding: utf-8 -*-
# Generated by Django 1.11.18 on 2019-02-18 05:24
from __future__ import unicode_literals

from django.db import migrations


class Migration(migrations.Migration):

    dependencies = [
        ('Lt', '0001_initial'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='category',
            name='admin',
        ),
        migrations.RemoveField(
            model_name='thumbup',
            name='article',
        ),
        migrations.RemoveField(
            model_name='thumbup',
            name='user',
        ),
        migrations.DeleteModel(
            name='Category',
        ),
        migrations.DeleteModel(
            name='ThumbUp',
        ),
    ]
