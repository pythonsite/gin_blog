package models

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
)

// query result
type QrArchive struct {
	ArchiveDate time.Time //month
	Total       int       //total
	Year        int       // year
	Month       int       // month
}

func MustListPostArchives() []*QrArchive {
	archives, _ := ListPostArchives()
	return archives
}

func ListPostArchives()([]*QrArchive, error) {
	var archives []*QrArchive
	querysql := `select DATE_FORMAT(created_at,'%Y-%m') as month, count(*) as total from posts where is_published=? group by month order by month desc`
	rows, err := DB.Raw(querysql, true).Rows()
	if err != nil {
		logs.Error("sql  %s error:%s", querysql, err)
		return  nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var archive QrArchive
		var month string
		_ = rows.Scan(&month, &archive.Total)
		//_ = DB.ScanRows(rows, &archive)
		archive.ArchiveDate, _ = time.Parse("2006-01",month)
		archive.Year = archive.ArchiveDate.Year()
		archive.Month = int(archive.ArchiveDate.Month())
		archives = append(archives, &archive)
	}
	logs.Info("all archives is:%#v", archives)
	return archives, nil
}

func ListPostByArchive(year, month string, pageIndex, pageSize int) ([]*Post,error) {
	var (
		rows *sql.Rows
		err error
	)
	if len(month) == 1 {
		month = "0" + month
	}
	condition := fmt.Sprintf("%s-%s",year, month)
	logs.Info(condition)
	if pageIndex > 0 {
		quersql := `select * from posts where date_format(created_at,'%Y-%m') = ? and is_published = ? order by created_at desc limit ? offset ?`
		rows, err = DB.Raw(quersql, condition, true, pageSize, (pageIndex-1)*pageSize).Rows()
	} else {
		querysql := `select * from posts where date_format(created_at,'%Y-%m') = ? and is_published = ? order by created_at desc`
		rows, err = DB.Raw(querysql, condition, true).Rows()
	}
	if err != nil {
		logs.Error("error:%v",err)
		return nil, err
	}
	defer rows.Close()
	posts := make([]*Post, 0)
	for rows.Next() {
		var post Post
		_ = DB.ScanRows(rows, &post)
		posts = append(posts, &post)
	}
	return posts, nil
}

func CountPostByArchive(year, month string) (count int, err error) {
	if len(month) == 1 {
		month = "0" + month
	}
	condition := fmt.Sprintf("%s-%s", year, month)
	quersql := `select count(*) from posts where date_format(created_at, '%Y-%m') = ? and is_published = ? order by created_at desc`
	err = DB.Raw(quersql, condition,true).Row().Scan(&count)
	return
}