package application

import (
	"context"
	"log"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"

	"github.com/andriykutsevol/DDDCasbinExample/configs"
	"github.com/andriykutsevol/DDDCasbinExample/internal/domain/rbac"
)

type AutoLoadPolicyChan chan *PolicyItem

type PolicyItem struct {
	Ctx      context.Context
	Enforcer *casbin.SyncedEnforcer
}

type RbacAdapter interface {
	persist.Adapter
	CreateAutoLoadPolicyChan() AutoLoadPolicyChan
	GetAutoLoadPolicyChan() AutoLoadPolicyChan
	AddPolicyItemToChan(ctx context.Context, e *casbin.SyncedEnforcer)
}

type rbacAdapter struct {
	rbacRepo rbac.Repository
}

func NewRbacAdapter(rbacRepo rbac.Repository) RbacAdapter {
	return &rbacAdapter{
		rbacRepo: rbacRepo,
	}
}

var autoLoadPolicyChan AutoLoadPolicyChan

func (a *rbacAdapter) CreateAutoLoadPolicyChan() AutoLoadPolicyChan {

	autoLoadPolicyChan = make(chan *PolicyItem, 1)

	go func() {
		for item := range autoLoadPolicyChan {
			err := item.Enforcer.LoadPolicy()
			if err != nil {
				//logger.WithContext(item.Ctx).Errorf("The load casbin policy error: %s", err.Error())
				log.Printf("The load casbin policy error: %s", err.Error())
			}
		}
	}()
	return autoLoadPolicyChan
}

func (a *rbacAdapter) GetAutoLoadPolicyChan() AutoLoadPolicyChan {
	return autoLoadPolicyChan
}

func (a *rbacAdapter) AddPolicyItemToChan(ctx context.Context, e *casbin.SyncedEnforcer) {
	if !configs.C.Casbin.Enable {
		return
	}

	if len(autoLoadPolicyChan) > 0 {
		//logger.WithContext(ctx).Infof("The load casbin policy is already in the wait queue")
		log.Printf("The load casbin policy is already in the wait queue")
		return
	}

	autoLoadPolicyChan <- &PolicyItem{
		Ctx:      ctx,
		Enforcer: e,
	}
}

func casbinModelLog(model *casbinModel.Model) {
	// var modelInfo [][]string
	// for k, v := range *model {
	// 	if k == "logger" {
	// 		continue
	// 	}

	// 	for i, j := range v {
	// 		modelInfo = append(modelInfo, []string{k, i, j.Value})
	// 	}
	// }
}

func casbinPolicyLog(model *casbinModel.Model) {
	policy := make(map[string][][]string)

	for key, ast := range (*model)["p"] {
		value, found := policy[key]
		if found {
			value = append(value, ast.Policy...)
			policy[key] = value
		} else {
			policy[key] = ast.Policy
		}
	}

	for key, ast := range (*model)["g"] {
		value, found := policy[key]
		if found {
			value = append(value, ast.Policy...)
			policy[key] = value
		} else {
			policy[key] = ast.Policy
		}
	}

}

func (a *rbacAdapter) LoadPolicy(model casbinModel.Model) error {
	casbinModelLog(&model)

	ctx := context.Background()

	policies, err := a.rbacRepo.ListRolesPolices(ctx)
	if err != nil {
		log.Printf("Load casbin role policy error: %s", err.Error())
		//logger.WithContext(ctx).Errorf("Load casbin role policy error: %s", err.Error())
		return err
	}
	for _, policy := range policies {
		err := persist.LoadPolicyLine(policy, model)
		if err != nil {
			log.Panic("ERROR: persist.LoadPolicyLine: ", err)
		}

	}
	casbinPolicyLog(&model)
	//map[g:[] p:[[menu_admin /api/v1/menus POST]]]

	policies, err = a.rbacRepo.ListUsersPolices(ctx)
	if err != nil {
		log.Printf("Load casbin user policy error: %s", err.Error())
		//logger.WithContext(ctx).Errorf("Load casbin user policy error: %s", err.Error())
		return err
	}
	for _, policy := range policies {
		err := persist.LoadPolicyLine(policy, model)
		if err != nil {
			log.Panic("ERROR: persist.LoadPolicyLine: ", err)
		}
	}
	casbinPolicyLog(&model)
	//map[g:[] p:[[menu_admin /api/v1/menus POST]]]

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *rbacAdapter) SavePolicy(model casbinModel.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *rbacAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *rbacAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *rbacAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
