package dal

import (
	"github.com/yzx9/otodo/model/dto"
	"github.com/yzx9/otodo/model/entity"
	"github.com/yzx9/otodo/util"
)

func InsertTodoList(todoList *entity.TodoList) error {
	re := db.Create(todoList)
	return util.WrapGormErr(re.Error, "todo list")
}

func SelectTodoList(id int64) (entity.TodoList, error) {
	var list entity.TodoList
	where := entity.TodoList{Entity: entity.Entity{ID: id}}
	re := db.Where(&where).First(&list)
	return list, util.WrapGormErr(re.Error, "todo list")
}

func SelectTodoLists(userId int64) ([]entity.TodoList, error) {
	var lists []entity.TodoList
	re := db.Where(entity.TodoList{UserID: userId}).Find(&lists)
	return lists, util.WrapGormErr(re.Error, "todo list")
}

func SelectTodoListsWithMenuFormat(userID int64) ([]dto.TodoListMenuItemRaw, error) {
	var lists []dto.TodoListMenuItemRaw
	re := db.
		Model(entity.TodoList{}).
		Where(entity.TodoList{UserID: userID}).
		Not(entity.TodoList{IsBasic: true}). // Skip basic todo list
		Select("id", "name", "todo_list_folder_id", "(SELECT count(todos.id) FROM todos WHERE todos.todo_list_id = todo_lists.id) as count").
		Find(&lists)
	return lists, util.WrapGormErr(re.Error, "todo list")
}

func SaveTodoList(todoList *entity.TodoList) error {
	re := db.Save(&todoList)
	return util.WrapGormErr(re.Error, "todo list")
}

func DeleteTodoList(id int64) error {
	re := db.Delete(&entity.Todo{Entity: entity.Entity{ID: id}})
	return util.WrapGormErr(re.Error, "todo list")
}

func DeleteTodoListsByFolder(todoListFolderID int64) (int64, error) {
	re := db.Where(entity.TodoList{TodoListFolderID: todoListFolderID}).Delete(entity.TodoList{})
	return re.RowsAffected, util.WrapGormErr(re.Error, "todo list")
}

func ExistTodoList(id int64) (bool, error) {
	var count int64
	where := entity.TodoList{Entity: entity.Entity{ID: id}}
	re := db.Model(&entity.TodoList{}).Where(&where).Count(&count)
	return count != 0, util.WrapGormErr(re.Error, "todo list")
}

/**
 * Sharing
 */

func InsertTodoListSharedUser(userID, todoListID int64) error {
	user := entity.User{Entity: entity.Entity{ID: userID}}
	list := entity.TodoList{Entity: entity.Entity{ID: todoListID}}
	err := db.Model(&user).Association("SharedTodoLists").Append(&list)
	return util.WrapGormErr(err, "todo list shared user")
}

func SelectSharedTodoLists(userID int64) ([]entity.TodoList, error) {
	user := entity.User{Entity: entity.Entity{ID: userID}}
	var lists []entity.TodoList
	err := db.Model(&user).Association("SharedTodoLists").Find(&lists)
	return lists, util.WrapGormErr(err, "user shared todo list")
}

func SelectTodoListSharedUsers(todoListID int64) ([]entity.User, error) {
	list := entity.TodoList{Entity: entity.Entity{ID: todoListID}}
	var users []entity.User
	err := db.Model(&list).Association("SharedUsers").Find(&users)
	return users, util.WrapGormErr(err, "todo list shared users")
}

func DeleteTodoListSharedUser(userID, todoListID int64) error {
	user := entity.User{Entity: entity.Entity{ID: userID}}
	list := entity.TodoList{Entity: entity.Entity{ID: todoListID}}
	err := db.Model(&list).Association("SharedUsers").Delete(&user)
	return util.WrapGormErr(err, "todo list shared users")
}

func ExistTodoListSharing(userID, todoListID int64) (bool, error) {
	// TODO[pref]: count in db
	user := entity.User{Entity: entity.Entity{ID: userID}}
	list := entity.TodoList{Entity: entity.Entity{ID: todoListID}}
	var lists []entity.TodoList
	if err := db.Model(&user).Association("SharedTodoLists").Find(&lists, &list); err != nil {
		return false, util.WrapGormErr(err, "todo list sharing")
	}

	return len(lists) != 0, nil
}
