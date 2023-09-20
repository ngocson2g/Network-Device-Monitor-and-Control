#ifndef APPCHAT_H
#define APPCHAT_H

#include <QMainWindow>

#include "functionacction.h"

QT_BEGIN_NAMESPACE
namespace Ui { class Appchat; }
QT_END_NAMESPACE

class Appchat : public QMainWindow
{
    Q_OBJECT

public:
    Appchat(QWidget *parent = nullptr);
    ~Appchat();

signals:
    void send_msg( QString msg);
    void connect_enemy(QString enemy);
    void start_app();

private slots:
    void on_pushButton_start_clicked();

    void on_pushButton_connect_clicked();

    void on_pushButton_send_msg_clicked();

    void on_pushButton_back_clicked();

    void on_write_comboBox_list_connection(QString enemy);
    void on_write_textEdit_show_msg(bool user, QString msg);

    void on_tester(QString str);

    void on_connect_next();

private:
    Ui::Appchat *ui;
    FunctionAcction *Acction = new FunctionAcction();
};
#endif // APPCHAT_H
