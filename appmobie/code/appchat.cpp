#include "appchat.h"
#include "ui_appchat.h"

Appchat::Appchat(QWidget *parent)
    : QMainWindow(parent)
    , ui(new Ui::Appchat)
{
    ui->setupUi(this);

    QPixmap pix(":/resouce/img/backgroud.jpg");
    //    ui->label_backgroud->setPixmap(pix);
    //ui->label_bg2->setPixmap(pix);
    //ui->label_bg3->setPixmap(pix);

    // Kết nối các signal và slot tương ứng
    connect(this, SIGNAL(start_app()), Acction, SLOT(on_start_app()));
    connect(this, SIGNAL(send_msg(QString)), Acction, SLOT(on_send_msg(QString)));
    connect(this, SIGNAL(connect_enemy(QString)), Acction, SLOT(on_connect_enemy(QString)));

    connect(Acction, SIGNAL(write_comboBox_list_connection(QString)), this, SLOT(on_write_comboBox_list_connection(QString)));
    connect(Acction, SIGNAL(write_textEdit_show_msg(bool,QString)), this, SLOT(on_write_textEdit_show_msg(bool,QString)));

    //    connect(Acction, SIGNAL(connect_next()), this, SLOT(on_connect_next()));

    connect(Acction, SIGNAL(tester(QString)), this, SLOT(on_tester(QString)));
}

Appchat::~Appchat()
{
    delete ui;
}

// Xử lý sự kiện khi nút "Start" được nhấn
void Appchat::on_pushButton_start_clicked()
{
    emit start_app();
    ui->stackedWidget_appChat->setCurrentIndex(1);
    ui->label_my_id->setText(Acction->myAddress);
}

// Xử lý sự kiện khi nút "Connect" được nhấn
void Appchat::on_pushButton_connect_clicked()
{
    QString address = ui->comboBox_list_connect->currentText();
    if (address.length()) {
        emit connect_enemy(address);
        ui->stackedWidget_appChat->setCurrentIndex(2);
    } else {
        // Xử lý trường hợp không có địa chỉ
    }
}

// Xử lý sự kiện khi nút "Send" được nhấn
void Appchat::on_pushButton_send_msg_clicked()
{
    QString msg = ui->lineEdit_send_msg->text();
    if (msg.length()) {
        emit send_msg(msg);
    } else {
        // Xử lý trường hợp không có tin nhắn
    }
    ui->lineEdit_send_msg->clear();
}

// Hàm ghi danh sách kết nối vào comboBox
void Appchat::on_write_comboBox_list_connection(QString enemy)
{
    ui->comboBox_list_connect->addItem(enemy);
}

// Hàm ghi tin nhắn vào textEdit để hiển thị
void Appchat::on_write_textEdit_show_msg(bool user, QString msg)
{
    if(user) {
        msg = "Me: " + msg;
    }else {
        msg = "Enemy: " + msg;
    }
    ui->textEdit_show_msg->append(msg);
}

// Hàm hiển thị thông báo kiểm tra (tester)
void Appchat::on_tester(QString str)
{
    QMessageBox::information(this, "tester", str);
}

// Hàm xử lý sự kiện khi kết nối hoàn thành
void Appchat::on_connect_next()
{
    QMessageBox::information(this, "Notifiaction", "Connect complete!");
}

// Xử lý sự kiện khi nút "Back" được nhấn
void Appchat::on_pushButton_back_clicked()
{
    ui->stackedWidget_appChat->setCurrentIndex(1);
}
