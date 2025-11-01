package repository

//go:generate sh -c "rm -rf mocks && mkdir -p mocks"
//go:generate minimock -i OrderRepository -o ./mocks/ -s "_minimock.go"
