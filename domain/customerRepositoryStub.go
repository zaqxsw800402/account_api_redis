package domain

type CustomerRepositoryStub struct {
	customers []Customer
}

func (s CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return s.customers, nil
}

//func NewCustomerRepositoryStub() CustomerRepositoryStub {
//	customers := []Customer{
//		{1000, "Lily", "Taiwan", "201", "2000-01-01", "1"},
//		//{"1002", "Ivy", "Taiwan", "201", "2000-01-01", "1"},
//	}
//	return CustomerRepositoryStub{customers}
//
//}
