# -*- coding:utf-8 -*-
__author__ = 'lanjie'
__Date__ = '2018/8/30 9:55'

from Base import BasePlayer,Chess,is_win
from PyQt5 import QtGui
from PyQt5.QtGui import QPixmap
from PyQt5.QtWidgets import QMessageBox,QLabel
# 生成逻辑棋盘


class DoublePlayer(BasePlayer):
    '''
    双人对战
    '''
    def __init__(self,parent=None):
        super().__init__(parent)
        # print("__init__")
        self.chessboard = [[None for i in range(0, 19)] for j in range(0, 19)]
        self.is_black = True
        self.is_over = False
        self.restart_btn.clicked.connect(self.restart)
        self.win_label = None
        # print(self.chessboard)

    def restart(self):
        # 重新开始游戏
        self.is_over = False
        # 清空胜利图片
        if self.win_label is not None:
            self.win_label.close()
        # 清空棋盘
        for i in range(0,19):
            for j in range(0,19):
                if self.chessboard[j][i] is not None:
                    self.chessboard[j][i].close()
                    self.chessboard[j][i] = None


    def mouseReleaseEvent(self, a0: QtGui.QMouseEvent):
        # 如果游戏已经结束，点击失效
        if self.is_over == True:
            return

        if a0.x() < 40 or a0.x() > 600:
            return
        if a0.y() < 40 or a0.y() > 600:
            return
        # 通过标识，决定棋子的颜色
        if self.is_black:
            self.chess = Chess(color='b',parent=self)
        else :
            self.chess = Chess('w',self)
        # 翻转棋子颜色
        self.is_black = not self.is_black
        # 将棋子定位到准确的坐标点
        if (a0.x() - 50 ) % 30 <= 15:
            x = (a0.x() - 50) // 30 * 30 + 50
        else:
            x = ((a0.x() - 50) // 30 + 1) * 30 + 50

        if (a0.y() - 50) % 30 <= 15:
            y = (a0.y() - 50) // 30 * 30 + 50
        else:
            y = ((a0.y() - 50) // 30 + 1) * 30 + 50
        # 在棋盘数组中，保存棋子对象
        xx = (x-50) // 30
        yy = (y-50) // 30
        # 如果此处已经有棋子，点击失效
        if self.chessboard[xx][yy] is not None :
            return

        self.chessboard[xx][yy] = self.chess

        x = x - self.chess.width()/2
        y = y - self.chess.height()/2

        self.chess.move(x,y)
        self.chess.show()

        color = is_win(self.chessboard)
        if color is False:
            return
        else :
            # QMessageBox.information(self,"消息","{}棋胜利".format(color))
            self.win_label = QLabel(self)
            if color == 'b' :
                pic = QPixmap("source/黑棋胜利.png")
            else:
                pic = QPixmap("source/白棋胜利.png")
            self.win_label.setPixmap(pic)
            self.win_label.move(100,100)
            self.win_label.show()
            self.is_over = True





