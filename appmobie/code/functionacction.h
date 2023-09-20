#ifndef FUNCTIONACCTION_H
#define FUNCTIONACCTION_H

#include <QObject>

#include <QtNetwork/QTcpServer>
#include <QtNetwork/QTcpSocket>
#include <QMessageBox>
#include <QList>
#include <QByteArray>
#include <QString>
#include <QHostAddress>
#include <QtNetwork/QNetworkInterface>
#include <QDebug>

#include <QPixmap>

class FunctionAcction : public QObject
{
    Q_OBJECT
public:
    explicit FunctionAcction(QObject *parent = nullptr);

signals:
    void write_comboBox_list_connection(QString enemy);
    void write_textEdit_show_msg(bool user, QString msg);

    void tester(QString str);
    void connect_next();

public slots:
    void on_send_msg(QString msg);
    void on_connect_enemy(QString enemy);
    void on_start_app();

    bool readData();
    void newConnection();
    void addConnection(QTcpSocket *socket);

    void get_Address(); //my address ...database

    void get_ListAddress_Server();

    void check_status_app();

    bool connect_confirm();

public:
    QTcpServer *Server = new QTcpServer;
    QTcpSocket *Client = new QTcpSocket;

    QList<QTcpSocket*> listClient;
    QList<QString> listServer;

    int Host = 8080; //default host setting

    bool app_status = true; // true: client ; false : server

    QString myAddress;
    QString enemyAddress;

    QString pass = "#$%#@#$";
};

#endif // FUNCTIONACCTION_H
