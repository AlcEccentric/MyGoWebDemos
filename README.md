# MyGoWebDemos
Some web demos extended from the hands-ons of the Udemy course GoWebDev

GoWebDemo1:
- Implement a backend which implements session and cookie
- Additional function: count how many times a user have visited this site
- Some notes from this demo:
	1. Reason of "http: multiple response.WriteHeader calls":
	See: https://stackoverflow.com/questions/27972715/multiple-response-writeheader-calls-in-really-simple-example
	2. Reason that the template cannot use .memberName to get access to the data in a struct
	Answer: the first letter of member name should be uppercase, which means the struct allows this member to be exported
	3. Reason that u.VisitCount++ cannot work
	Answer: The dbUsers map does not use reference to user as its value. So when getting user from it and modifying content of that intance, I am actually changing data of a copy of that user not the original one in the map	
