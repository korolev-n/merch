package repository

type UserRepository interface {
    Create() (error)
}

type WalletRepository interface {
    Create() error
}