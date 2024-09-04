package menu

import (
	"context"
	"fmt"
	"log"
	"strings"

	//"github.com/google/uuid"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/menu"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/pagination"
	"github.com/andriykutsevol/DDDCasbinExample/internal/infrastructure/pg/storage"
)

func NewRepository(service storage.DatabaseService) menu.Repository {
	return &repository{
		service: service,
	}
}

type repository struct {
	service storage.DatabaseService
}

func (r *repository) Query(ctx context.Context, params menu.QueryParam) (menu.Menus, *pagination.Pagination, error) {
	// Select all
	// When you execute a query like SELECT * FROM users using pgx,
	// the columns in each row are returned in the same order as they are defined in the table schema.
	// This order is determined by the sequence in which the columns were originally created in the table.

	var whereClauses []string
	var queryParams []interface{}
	var paramIdx int = 1

	if v := params.IDs; len(v) > 0 {
		//db = db.Where("id IN (?)", v)

		// Convert the slice to an interface slice for pgx
		args := make([]interface{}, len(params.IDs))
		for i, v := range params.IDs {
			args[i] = v
		}
		whereClause := fmt.Sprintf("id = ANY($%d)", paramIdx)
		whereClauses = append(whereClauses, whereClause)
		queryParams = append(queryParams, args)
		paramIdx++
	}
	if v := params.Name; v != "" {
		//db = db.Where("name=?", v)

		whereClause := fmt.Sprintf("name = $%d", paramIdx)
		whereClauses = append(whereClauses, whereClause)
		queryParams = append(queryParams, params.Name)
		paramIdx++
	}
	if v := params.ParentID; v != "" {
		//db = db.Where("parent_id=?", *v)

		whereClause := fmt.Sprintf("parentid = $%d", paramIdx)
		whereClauses = append(whereClauses, whereClause)
		queryParams = append(queryParams, params.Name)
		paramIdx++
	}
	if v := params.PrefixParentPath; v != "" {
		//db = db.Where("parent_path LIKE ?", v+"%")

		whereClause := fmt.Sprintf("parentpath LIKE $%d", paramIdx)
		whereClauses = append(whereClauses, whereClause)
		queryParams = append(queryParams, params.Name)
		paramIdx++
	}
	if v := params.ShowStatus; v != 0 {
		//db = db.Where("show_status=?", v)

		whereClause := fmt.Sprintf("showstatus = $%d", paramIdx)
		whereClauses = append(whereClauses, whereClause)
		queryParams = append(queryParams, params.Name)
		paramIdx++

	}
	if v := params.Status; v != 0 {
		//db = db.Where("status=?", v)

		whereClause := fmt.Sprintf("status = $%d", paramIdx)
		whereClauses = append(whereClauses, whereClause)
		queryParams = append(queryParams, params.Name)
		paramIdx++
	}
	if v := params.QueryValue; v != "" {
		// v = "%" + v + "%"
		// db = db.Where("name LIKE ? OR memo LIKE ?", v, v)

		//v = "%" + v + "%"
		whereClause := fmt.Sprintf("name LIKE $%d OR memo LIKE $%d", paramIdx, paramIdx+1)
		whereClauses = append(whereClauses, whereClause)
		queryParams = append(queryParams, params.Name)
		// paramIdx += 2
		// _ = paramIdx		// TODO: it is for liter.
	}

	queryString := "SELECT * FROM casbin.menu"
	if len(whereClauses) > 0 {
		queryString = "SELECT * FROM casbin.menu WHERE "
		queryString += strings.Join(whereClauses, " AND ")
	}

	rows, err := r.service.Pool.Query(context.Background(), queryString, queryParams...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var menus []*Model

	for rows.Next() {
		var menu Model
		err := rows.Scan(
			&menu.ID,
			&menu.Name,
			&menu.Sequence,
			&menu.Icon,
			&menu.Router,
			&menu.ParentID,
			&menu.ParentPath,
			&menu.ShowStatus,
			&menu.Status,
			&menu.Memo,
			&menu.Creator,
			&menu.CreatedAt,
			&menu.UpdatedAt,
			&menu.DeletedAt,
			&menu.IDString,
		)
		if err != nil {
			return nil, nil, err
		}
		menus = append(menus, &menu)
	}
	return toDomainList(menus), nil, nil
}

func (r *repository) Get(ctx context.Context, id string) (*menu.Menu, error) {
	var menu Model
	queryString := "SELECT * FROM casbin.menu WHERE id = $1"
	err := r.service.Pool.QueryRow(context.Background(), queryString, id).Scan(
		&menu.ID,
		&menu.Name,
		&menu.Sequence,
		&menu.Icon,
		&menu.Router,
		&menu.ParentID,
		&menu.ParentPath,
		&menu.ShowStatus,
		&menu.Status,
		&menu.Memo,
		&menu.Creator,
		&menu.CreatedAt,
		&menu.UpdatedAt,
		&menu.DeletedAt,
		&menu.IDString,
	)
	if err != nil {
		return nil, err
	}
	return menu.ToDomain(), nil
}

func (r *repository) GetByIdString(ctx context.Context, id string) (*menu.Menu, error) {
	var menu Model
	queryString := "SELECT * FROM casbin.menu WHERE idstring = $1"
	err := r.service.Pool.QueryRow(context.Background(), queryString, id).Scan(
		&menu.ID,
		&menu.Name,
		&menu.Sequence,
		&menu.Icon,
		&menu.Router,
		&menu.ParentID,
		&menu.ParentPath,
		&menu.ShowStatus,
		&menu.Status,
		&menu.Memo,
		&menu.Creator,
		&menu.CreatedAt,
		&menu.UpdatedAt,
		&menu.DeletedAt,
		&menu.IDString,
	)
	if err != nil {
		return nil, err
	}
	return menu.ToDomain(), nil
}

func (r *repository) Create(ctx context.Context, item *menu.Menu) error {

	model := domainToModel(item)
	queryString := `
		INSERT INTO casbin.menu (
		id, 
		name, 
		sequence, 
		icon, 
		router, 
		parentid, 
		parentpath, 
		showstatus, 
		status, 
		memo, 
		creator, 
		createdat, 
		updatedat,
		deletedat, 
		idstring
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	_, err := r.service.Pool.Exec(context.Background(), queryString,
		model.ID,
		model.Name,
		model.Sequence,
		model.Icon,
		model.Router,
		model.ParentID,
		model.ParentPath,
		model.ShowStatus,
		model.Status,
		model.Memo,
		model.Creator,
		model.CreatedAt,
		model.UpdatedAt,
		model.DeletedAt,
		model.IDString)
	if err != nil {
		log.Fatalf("Failed to insert data: %v\n", err)
	}
	return nil
}

func (r *repository) Update(ctx context.Context, id string, item *menu.Menu) error {
	log.Println("menu/repository.go: Update")
	return nil
}

func (r *repository) UpdateParentPath(ctx context.Context, id, parentPath string) error {
	log.Println("menu/repository.go: UpdateParentPath")
	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	log.Println("menu/repository.go: Delete")
	return nil
}

func (r *repository) UpdateStatus(ctx context.Context, id string, status int) error {
	log.Println("menu/repository.go: UpdateStatus")
	return nil
}

func (r *repository) Purge(ctx context.Context) error {
	_, err := r.service.Pool.Exec(context.Background(), "DELETE FROM casbin.menu")
	if err != nil {
		log.Printf("Failed to delete from table %s: %v\n", "menu", err)
	}

	return nil
}
