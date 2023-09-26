package models

import (
	"database/sql"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type ExampleModel struct {
	DB *sql.DB
}

func (m *ExampleModel) ExampleTransaction() error {
	// Вызов метода Begin() в пуле соединений создает новый объект sql.Tx
	// который представляет текущую транзакцию к базы данных.
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	// Вызываем Exec() для транзакции, передавая оператор и любые другие
	// параметры. Важно отметить, что tx.Exec() вызывается для
	// только что созданного объекта транзакции, а НЕ для пула соединений. Хотя мы
	// здесь используем tx.Exec(), вы также можете использовать tx.Query() и tx.QueryRow()
	// таким же образом
	_, err = tx.Exec("INSERT INTO ...")
	if err != nil {
		// Если возникла ошибка, вызываем метод tx.Rollback() для
		// транзакции. Он прервет транзакцию и
		// в базу данных не будут внесены изменения.
		tx.Rollback()
		return err
	}

	// Точно так же выполняется другая транзакция.
	_, err = tx.Exec("UPDATE ...")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Если ошибок нет, то запрос транзакции может быть
	// выполнен в базе данных с помощью метода tx.Commit(). Очень важно ВСЕГДА
	// вызывать Rollback() или Commit() в конце функции перед "return". Если вы
	// этого не сделаете, соединение останется открытым и не будет возвращено
	// в пул соединений. Это может привести к достижению максимального лимита соединений и исчерпанию ресурсов.
	err = tx.Commit()
	return err
}
