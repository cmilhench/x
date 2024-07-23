package pg

/*

type DB struct {
	*sqlx.DB
}

func NewDB(driverName, dataSourceName string) (*DB, error) {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{
		db,
	}, nil
}

func NamedSelect[T any](ctx context.Context, d *DB, query string, arg any, out T) error {
	if arg == nil {
		arg = struct{}{}
	}
	value := reflect.ValueOf(out)
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	err = d.SelectContext(ctx, out, d.Rebind(query), args...)
	return err
}

*/
