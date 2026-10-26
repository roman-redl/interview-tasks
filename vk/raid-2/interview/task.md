# VK — заход 2, техническое интервью

- Даты: 2026-05-04
- Формат: онлайн, теория + разбор кода на экране

## 1. Ревью кода: перевод между аккаунтами

На экране показан код на Go. Есть система, в которой можно совершать переводы между аккаунтами: абстракция аккаунта (идентификатор, баланс, мьютекс — для «транзакционности») и абстракция перевода (от кого, кому и сумма), плюс функция выполнения перевода. Требуется оценить реализацию: увидеть ли проблемы, и найденное обсудить.

Код (файл `tasks/transfer.go`, package и импорты скрыты интерфейсом звонка):

```go
type Account struct {
	ID      int
	Balance uint
	MU      *sync.Mutex
}

type Transfer struct {
	From   *Account
	To     *Account
	Amount uint
}

func ExecuteTransfer(transfer Transfer) error {
	transfer.From.MU.Lock()
	defer transfer.From.MU.Unlock()

	transfer.To.MU.Lock()
	defer transfer.To.MU.Unlock()

	transfer.From.Balance -= transfer.Amount
	transfer.To.Balance += transfer.Amount

	return nil
}
```

## 2. Ревью кода: уникальные имена пользователей по проекту

Джун выполнил задачу и просит провести ревью кода. Его задание — написать две функции: первая ходит в PostgreSQL (через заглушку — копию интерфейса pgx) и достаёт все имена пользователей с фильтром по проекту; вторая убирает дубликаты имён из результата первой. В коде заведомо есть проблемы разного рода — от паник до странного кодстайла. Требуется найти все возможные проблемы и обсудить их.

Код (файл `tasks/first_names.go`, начало файла и хвост функции скрыты интерфейсом звонка):

```go
	"fmt"
)

type Conn interface {
	Query(ctx context.Context, q string, args ...any) (Rows, error)
}

type Rows interface {
	Next() bool
	Scan(destinations ...any) error
	Close() error
}

var conn Conn

func GetUniqueUserFirstNamesByProjectName(ctx context.Context, projectName string) ([]string, error) {
	userFirstNames, err := GetUserFirstNamesByProjectName(ctx, projectName)
	if err != nil {
		return nil, err
	}

	var (
		userFirstNamesMap    map[string]bool
		uniqueUserFirstNames []string
	)
	for _, userFirstName := range userFirstNames {
		if _, ok := userFirstNamesMap[userFirstName]; !ok {
			userFirstNamesMap[userFirstName] = true
			uniqueUserFirstNames = append(uniqueUserFirstNames, userFirstName)
		}
	// …
```
