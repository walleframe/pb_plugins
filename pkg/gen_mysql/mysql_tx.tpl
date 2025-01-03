// Code generated by {{.ToolName}} {{.Version}}. DO NOT EDIT.{{$tbl := .}}
package {{$tbl.DB}}

$Import-Packages$


var (
	global{{Title $tbl.DB}}DB atomic.Pointer[sqlx.DB]
)

func init() {
	{{$tbl.SvcDB}}.RegisterDB("mysql", "{{$tbl.DB}}", "_", func(db *sqlx.DB) error {
		global{{Title $tbl.DB}}DB.Store(db)
		return nil
	})
}

type {{Title $tbl.DB}}Tx struct {
    tx *sqlx.Tx
    db *sqlx.DB
}

func (tx *{{Title $tbl.DB}}Tx) Tx() *sqlx.Tx {
    return tx.tx
}

// Beginx begins a transaction and returns an *sqlx.Tx instead of an *sql.Tx.
func Begin() (_ *{{Title $tbl.DB}}Tx, err error) {
    db := global{{Title $tbl.DB}}DB.Load()
    tx,err := db.Beginx()
    if err != nil {
        return nil, err
    }
    return &{{Title $tbl.DB}}Tx{tx, db}, nil
}

// BeginTx begins a transaction and returns an *sqlx.Tx instead of an
// *sql.Tx.
//
// The provided context is used until the transaction is committed or rolled
// back. If the context is canceled, the sql package will roll back the
// transaction. Tx.Commit will return an error if the context provided to
// BeginxContext is canceled.
func BeginTx(ctx context.Context, opts *sql.TxOptions) (*{{Title $tbl.DB}}Tx, error) {
    db := global{{Title $tbl.DB}}DB.Load()
    tx,err := db.BeginTxx(ctx, opts)
    if err != nil {
        return nil,err
    }
    return &{{Title $tbl.DB}}Tx{tx, db}, nil
}


func PerformTx(ctx context.Context, handler func(ctx context.Context,tx *{{Title $tbl.DB}}Tx) error) error {
    tx,err := BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("{{$tbl.DB}}.PerformTx begin failed,%w", err)
    }
    defer func() {
        if r := recover(); r != nil {
            tx.Tx().Rollback()
            err = fmt.Errorf("{{$tbl.DB}}.PerformTx panic %v", r)
            return
        }
    }()
    err = handler(ctx, tx)
    if err != nil {
        tx.Tx().Rollback()
        return fmt.Errorf("{{$tbl.DB}}.PerformTx handle failed,%w", err)
    }
    err = tx.Tx().Commit()
    if err != nil {
        return fmt.Errorf("{{$tbl.DB}}.PerformTx commit failed,%w", err)
    }
    return nil
}