package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Transform struct {
	ID string `json:"id"`
	Name string `json:"name"`
	InputFormat string `json:"input_format"`
	OutputFormat string `json:"output_format"`
	Template string `json:"template"`
	Description string `json:"description"`
	RunCount int `json:"run_count"`
	LastRunAt string `json:"last_run_at"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"smelter.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS transforms(id TEXT PRIMARY KEY,name TEXT NOT NULL,input_format TEXT DEFAULT 'json',output_format TEXT DEFAULT 'json',template TEXT DEFAULT '',description TEXT DEFAULT '',run_count INTEGER DEFAULT 0,last_run_at TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Transform)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO transforms(id,name,input_format,output_format,template,description,run_count,last_run_at,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.InputFormat,e.OutputFormat,e.Template,e.Description,e.RunCount,e.LastRunAt,e.CreatedAt);return err}
func(d *DB)Get(id string)*Transform{var e Transform;if d.db.QueryRow(`SELECT id,name,input_format,output_format,template,description,run_count,last_run_at,created_at FROM transforms WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.InputFormat,&e.OutputFormat,&e.Template,&e.Description,&e.RunCount,&e.LastRunAt,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Transform{rows,_:=d.db.Query(`SELECT id,name,input_format,output_format,template,description,run_count,last_run_at,created_at FROM transforms ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Transform;for rows.Next(){var e Transform;rows.Scan(&e.ID,&e.Name,&e.InputFormat,&e.OutputFormat,&e.Template,&e.Description,&e.RunCount,&e.LastRunAt,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Transform)error{_,err:=d.db.Exec(`UPDATE transforms SET name=?,input_format=?,output_format=?,template=?,description=?,run_count=?,last_run_at=? WHERE id=?`,e.Name,e.InputFormat,e.OutputFormat,e.Template,e.Description,e.RunCount,e.LastRunAt,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM transforms WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM transforms`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Transform{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ? OR description LIKE ?)"
        args=append(args,"%"+q+"%");args=append(args,"%"+q+"%");
    }
    rows,_:=d.db.Query(`SELECT id,name,input_format,output_format,template,description,run_count,last_run_at,created_at FROM transforms WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Transform;for rows.Next(){var e Transform;rows.Scan(&e.ID,&e.Name,&e.InputFormat,&e.OutputFormat,&e.Template,&e.Description,&e.RunCount,&e.LastRunAt,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    return m
}
