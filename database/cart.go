package database

import (
	"errors"
)

var (
	ErrCantFindProduct    = errors.New("can't find product")
	ErrCantDecodeProducts = errors.New("can't find product")
	ErrUserIdIsNotValid   = errors.New("user id is not valid")
	ErrCantUpdateUser     = errors.New("can't update user")
	ErrCantRemoveItem     = errors.New("can't remove this item")
	ErrrCantGetItem       = errors.New("can't get item")
	ErrCantBuyCartItem    = errors.New("can't update the purchase")
)

func AddProductToCart() {

}
func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
