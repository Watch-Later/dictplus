package database

import (
	"database/sql"
	"fmt"

	"ahui2016.github.com/dictplus/model"
	"ahui2016.github.com/dictplus/stmt"
	"ahui2016.github.com/dictplus/util"

	_ "github.com/mattn/go-sqlite3"
)

type (
	Word = model.Word
)

type DB struct {
	Path string
	DB   *sql.DB
}

func (db *DB) mustBegin() *sql.Tx {
	tx, err := db.DB.Begin()
	util.Panic(err)
	return tx
}

func (db *DB) Exec(query string, args ...interface{}) (err error) {
	_, err = db.DB.Exec(query, args...)
	return
}

func (db *DB) Open(dbPath string) (err error) {
	if db.DB, err = sql.Open("sqlite3", dbPath+"?_fk=1"); err != nil {
		return
	}
	db.Path = dbPath
	if err = db.Exec(stmt.CreateTables); err != nil {
		return
	}
	return initFirstID(word_id_key, word_id_prefix, db.DB)
}

func (db *DB) GetWordByID(id string) (w Word, err error) {
	row := db.DB.QueryRow(stmt.GetWordByID, id)
	w, err = scanWord(row)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("not found (id:%s)", id)
	}
	return
}

func (db *DB) CountAllWords() (int64, error) {
	return getInt1(db.DB, stmt.CountAllWords)
}

func (db *DB) DeleteWord(id string) error {
	return db.Exec(stmt.DeleteWord, id)
}

func (db *DB) InsertNewWord(w *Word) (err error) {
	tx := db.mustBegin()
	defer tx.Rollback()

	if w.ID, err = getNextID(tx, word_id_key); err != nil {
		return
	}
	w.CTime = util.TimeNow()
	if err = insertWord(tx, w); err != nil {
		return
	}
	return tx.Commit()
}

func (db *DB) UpdateWord(w *Word) error {
	return updateWord(db.DB, w)
}

func (db *DB) GetWords(label, pattern string) (words []*Word, err error) {
	if label+pattern == "" {
		return nil, fmt.Errorf("nothing to search")
	}

	var rows *sql.Rows

	if label == "Recently-Added" {
		rows, err = db.DB.Query(stmt.NewWords, 30)
	}
	defer rows.Close()

	for rows.Next() {
		w, err := scanWord(rows)
		if err != nil {
			return nil, err
		}
		words = append(words, &w)
	}
	err = rows.Err()
	return
}
