# Be Couple #

A MVP project developed with Go. The purpose is to enhance my CV

### Components ###

* Web app
* Mobile app
* Backend

### generate database models ###
* database name: app_mvp_dating
* tool: https://github.com/xo/xo
* command:  <project_dir>$ xo mysql://root:qweasdzxc@123@localhost/app_mvp_dating -o ./models/xodb -f _xo.go --ignore-fields create_time update_time --template-path ./database/xo_templates/

### common mistakes
* MySQL database: if SELECT statement has BIT(1) field being converted to boolean type, add '+0'

### Daily basis
## Mac OS
# Open database connection
* enter command $ mysql.server start
* Open MySQL Bench, connect localhost
# Build server
* enter $ go build
# test
* enter for example: $ go test -run ^TestApiRegisterNew$
