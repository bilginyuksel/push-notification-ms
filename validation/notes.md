1. If you write a validator for a struct kind, it will effect all of the structs at every depth. 
    But also you can write special validators for particular structs like time. If you want to write a validator for time then
    you should add time.Time as a new Kind to map (it is already exist) then you should define the validator.