#include "functionacction.h"

FunctionAcction::FunctionAcction(QObject *parent)
    : QObject{parent}
{

}

// Xử lý khi gửi tin nhắn
void FunctionAcction::on_send_msg(QString msg)
{
    if (app_status) {
        if (Client->isOpen()) {
            Client->write(msg.toStdString().c_str());
            emit write_textEdit_show_msg(true, msg);
        }
    } else {
        foreach (QTcpSocket *socket, listClient) {
            quint32 addSocket = socket->peerAddress().toIPv4Address();

            if (QHostAddress(addSocket).toString() == enemyAddress) {
                socket->write((msg.toStdString().c_str()));
                emit write_textEdit_show_msg(true, msg);
            }
        }
    }
}

// Xử lý khi kết nối với đối thủ
void FunctionAcction::on_connect_enemy(QString enemy)
{
    enemyAddress = enemy;

    if (app_status) {
        Client->connectToHost(QHostAddress(enemyAddress), Host);
        connect(Client, SIGNAL(readyRead()), this, SLOT(readData()));
        Client->open(QIODevice::ReadWrite);

        // Xác nhận bắt tay (có thể bổ sung sau)
        // if (connect_confirm()) {
        //     emit connect_next();
        // }
    } /*else {
        if (readData()) {
            connect_next();
        }
    }*/
}

// Xử lý khi khởi động ứng dụng
void FunctionAcction::on_start_app()
{
    if (app_status) {
        check_status_app();
    }

    if (!app_status) {
        if(Server->listen(QHostAddress::Any, Host)) {
            connect(Server, SIGNAL(newConnection()), this, SLOT(newConnection()));
        }
    }

    if (app_status) { // Kiểm tra trạng thái ứng dụng (test)
        emit tester("client");
    } else {
        emit tester("server");
    }
}

// Xử lý khi nhận dữ liệu từ kết nối (Client hoặc Server)
bool FunctionAcction::readData()
{
    if (app_status) { // Client
        if (Client->isOpen()) {
            QByteArray Data = Client->readAll();
            QString msg = QString::fromStdString(Data.toStdString());

            if (msg.length()) {
                emit write_textEdit_show_msg(false, msg);
            }
        }
    } else { // Server
        QTcpSocket *socket = reinterpret_cast<QTcpSocket*> (sender());
        QByteArray Data = socket->readAll();
        QString msg = QString::fromStdString(Data.toStdString());

        if (msg.length()) {
            emit write_textEdit_show_msg(false, msg);
        }
    }
    return false;
}

// Xử lý khi có kết nối mới
void FunctionAcction::newConnection(){
    if (Server->hasPendingConnections()) {
        addConnection(Server->nextPendingConnection());
    }
}

// Thêm kết nối vào danh sách
void FunctionAcction::addConnection(QTcpSocket *socket)
{
    listClient.append(socket);
    connect(socket, SIGNAL(readyRead()), this, SLOT(readData()));
    quint32 address = socket->peerAddress().toIPv4Address();

    emit write_comboBox_list_connection(QHostAddress(address).toString());
}

// Lấy địa chỉ IP của máy
void FunctionAcction::get_Address()
{
    QList<QHostAddress> ipAddressList = QNetworkInterface::allAddresses();
    foreach (const QHostAddress& address, ipAddressList) {
        if (address != QHostAddress::LocalHost && address.toIPv4Address()) {
            myAddress = address.toString();
        }
    }
}

// Lấy danh sách địa chỉ IP của các Server
void FunctionAcction::get_ListAddress_Server()
{
    QString Laddress = myAddress.left(myAddress.lastIndexOf(".") + 1);

    Laddress = "127.0.0."; // Test

    QTcpSocket *socket = new QTcpSocket();
    for (int i = 1; i < 2; i++) {

        QString Address = Laddress + QString::number(i);

        if (Address == myAddress) {
            continue;
        }

        socket->connectToHost(QHostAddress(Address), Host);

        if (socket->waitForConnected(100)) {
            listServer.append(Address);
        }

        socket->disconnectFromHost();
    }
}

// Kiểm tra trạng thái ứng dụng
void FunctionAcction::check_status_app()
{
    get_Address();
    get_ListAddress_Server();

    if (!listServer.size()) {
        app_status = false;
    } else {
        foreach (QString address, listServer) {
            emit write_comboBox_list_connection(address);
        }
    }
}

// Xác nhận kết nối (có thể bổ sung sau)
bool FunctionAcction::connect_confirm()
{
    Client->write(pass.toStdString().c_str());
    if (readData()) {
        return true;
    }
    return false;
}
