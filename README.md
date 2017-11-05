# Be Couple #

A MVP project developed with Go. The purpose is to enhance my CV

### Components ###

* Web app
* Mobile app
* Backend

### api specification ###

## register : /api/register
# request
- primaryID (represent unique identifier, such as email)
- password
- fullname
# response
- on success, return normal response

   After registration, the use must confirm their account, using confirm token received via their email

## confirm : /api/confirm
# request 
- JWT (Header 'Bearer')
- __email__ (generated from JWT)
- confirm\_token
# response
- on success, return normal response

   After this step, the user can login their account

## authenticate : /api/auth
# request
- primaryID
- password
# response

   {
   	success: 1,
	data: {
	    token: {JSon Web Token here}
	}
   }
