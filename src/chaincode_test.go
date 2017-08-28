package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var pedidoID = "la1"
var pedidoJson = `{"cpf": "09596397729", "DescricaoItens": "Maquina Lavar Brastemp; Panela Tramontina", "ItensId": "234;445", "dataVenda": 1503849607000 }`


func ObterPedidoForTest( t *testing.T, stub shim.ChaincodeStubInterface, id string, p *Pedido){
	bytes, err := stub.GetState(id)
	if err != nil {
		t.Fatalf("Could not fetch pedido with ID " + pedidoID)
	}
	err = json.Unmarshal(bytes, &p)
	if err != nil {
		t.Fatalf("Could not unmarshal pedido with ID" + pedidoID)
	}
}



func TestCriarChaincode(t *testing.T) {
	fmt.Println("Entering TestCreateLoanApplication")
	attributes := make(map[string][]byte)
	//Create a custom MockStub that internally uses shim.MockStub
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}
}

func TestRegistrarPedidoErro(t *testing.T) {
	fmt.Println("Entering TestRegistrarPedidoErro")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	stub.MockTransactionStart("t123")
	_, err := RegistrarPedido(stub, []string{})
	if err == nil {
		t.Fatalf("Expected RegistrarPedido to return validation error")
	}
	stub.MockTransactionEnd("t123")
}

func TestRegistrarPedidoErroFormato(t *testing.T) {
	fmt.Println("Entering TestRegistrarPedidoErro")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	stub.MockTransactionStart("t123")
	_, err := RegistrarPedido(stub, []string{pedidoID, " blabla bla "})
	if err == nil {
		t.Fatalf("Expected RegistrarPedido to return validation error")
	}
	stub.MockTransactionEnd("t123")
}

func TestInvokePedidoSucesso(t *testing.T) {
	fmt.Println("Entering TestRegistrarPedidoSucesso")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	bytes, err := stub.MockInvoke("t123", "RegistrarPedido", []string{pedidoID, pedidoJson})
	if err != nil {
		t.Fatalf("Expected RegistrarPedido function to be invoked")
	}

	var pe Pedido
	bytes, err = stub.GetState(pedidoID)
	if err != nil {
		t.Fatalf("Could not fetch pedido with ID " + pedidoID)
	}
	err = json.Unmarshal(bytes, &pe)
	if err != nil {
		t.Fatalf("Could not unmarshal pedido with ID" + pedidoID)
	}

	var errors = []string{}
	var pedidoTestData Pedido
	err = json.Unmarshal([]byte(pedidoJson), &pedidoTestData)

	if pe.ID != pedidoID {
		errors = append(errors, "Pedido ID does not match")
	}
	if pe.CPFCliente != pedidoTestData.CPFCliente {
		errors = append(errors, "Pedido CPFCliente does not match")
	}
	if pe.DescricaoItens != pedidoTestData.DescricaoItens {
		errors = append(errors, "Pedido DescricaoItens does not match")
	}

	//Can be extended for all fields
	if len(errors) > 0 {
		t.Fatalf("Mismatch between input and stored Pedido")
		for j := 0; j < len(errors); j++ {
			fmt.Println(errors[j])
		}
	}

}

func TestUpdatePedidoSucesso(t *testing.T) {
	fmt.Println("Entering TestUpdatePedidoSucesso")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	_, err := stub.MockInvoke("t123", "RegistrarPedido", []string{pedidoID, pedidoJson})
	if err != nil {
		t.Fatalf("Expected RegistrarPedido function to be invoked")
	}
	
	novoCpf := "444";
	fn := func(p *Pedido) error {
		p.CPFCliente = novoCpf;
		return nil
	}
 	
 	stub.MockTransactionStart("t123")
	AtualizarPedido(stub, pedidoID, fn)
  	stub.MockTransactionEnd("t123")

	var pe Pedido
	ObterPedidoForTest(t, stub, pedidoID, &pe);	
	if pe.CPFCliente != novoCpf {
		t.Fatalf("CPF not updated")
	}
}

func TestInvokeValidationError(t *testing.T) {
	fmt.Println("Entering TestInvokeValidation")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	_, err := stub.MockInvoke("t123", "BlaBlaBla", []string{})
	if err == nil {
		t.Fatalf("Expected unknow invoke method")
	}

}

func TestInvokeValidationSuccess(t *testing.T) {
	fmt.Println("Entering TestInvokeValidation")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	bytes, err := stub.MockInvoke("t123", "RegistrarPedido", []string{pedidoID, pedidoJson})
	if err != nil {
		t.Fatalf("Expected RegistrarPedido function to be invoked")
	}

	var pe Pedido
	err = json.Unmarshal(bytes, &pe)
	if err != nil {
		t.Fatalf("Expected valid Pedido JSOn string to be returned from RegistrarPedido method")
	}

}


func TestQueryErro(t *testing.T) {
		fmt.Println("Entering TestInvokeValidation")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	_, err := stub.MockQuery("BlaBlaBla", []string{})
	if err == nil {
		t.Fatalf("Expected unknow query method")
	}
}



func TestQuerySucesso(t *testing.T) {
	fmt.Println("Entering TestQueryErro")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	_, err := stub.MockInvoke("t123", "RegistrarPedido", []string{pedidoID, pedidoJson})
	if err != nil {
		t.Fatalf("Expected RegistrarPedido function to be invoked")
	}

	bytes, err := stub.MockQuery("ObterPedido", []string{pedidoID})
	if err != nil {
		t.Fatalf("Expected ObterPedido function to be invoked correctly")
	}

	var pe Pedido
	err = json.Unmarshal(bytes, &pe)
	if err != nil {
		t.Fatalf("Could not unmarshal pedido with ID" + pedidoID)
	}
	if pe.ID != pedidoID {
		t.Fatalf("Not query successfully: " + pedidoID)
	}
}

func TestRegistrarEntregaErro(t * testing.T) {
	fmt.Println("Entering TestRegistrarEntrega")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}

	_, err := stub.MockInvoke("t123", "RegistrarEntrega", []string{pedidoID})
	if err == nil {
		t.Fatalf("Expected TestRegistrarEntrega give a error")
	}


}


func TestRegistrarEntregaErroData(t * testing.T) {
	fmt.Println("Entering TestRegistrarEntregaErroData")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}
	_, err := stub.MockInvoke("t123", "RegistrarPedido", []string{pedidoID, pedidoJson})

	_, err = stub.MockInvoke("t123", "RegistrarEntrega", []string{pedidoID, "eeeeeee"})
	if err == nil {
		t.Fatalf("Expected TestRegistrarEntrega give a error")
	}	
}

func TestRegistrarEntregaErroSucesso(t * testing.T) {
	fmt.Println("Entering TestRegistrarEntregaErroSucesso")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)
	if stub == nil {
		t.Fatalf("MockStub creation failed")
	}
	
	stub.MockInvoke("t123", "RegistrarPedido", []string{pedidoID, pedidoJson})

	stub.MockInvoke("t123", "RegistrarEntrega", []string{pedidoID, "654"})
	
	var pe Pedido
	ObterPedidoForTest(t, stub, pedidoID, &pe);
	if pe.DataEntrega != 654 {
		t.Fatalf("Data Entrega not updated")
	}
}

func TestArrependimentoErro(t * testing.T) {
	fmt.Println("Entering TestRegistrarEntregaErroSucesso")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)

	_, err := stub.MockInvoke("t123", "Arrependimento", []string{pedidoID})

	if err == nil {
		t.Fatalf("Expected error in call Arrependimento")
	}
}

func TestArrependimentoSemRegistroEntrega(t * testing.T) {
	fmt.Println("Entering TestArrependimentoErroDataAntes7Dias")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)

	stub.MockInvoke("t123", "RegistrarPedido", []string{pedidoID, pedidoJson})
	_, err := stub.MockInvoke("t123", "RegistrarArrependimento", []string{pedidoID, "1503849607000"})
	if err == nil {
		t.Fatalf("Expected not delivery error ")
	}

}

func TestArrependimentoErroDataDepois7Dias(t * testing.T) {
	fmt.Println("Entering TestArrependimentoErroDataAntes7Dias")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)

	stub.MockInvoke("t123", "RegistrarPedido", []string{pedidoID, pedidoJson})

	stub.MockInvoke("t123", "RegistrarEntrega", []string{pedidoID, "1472313607000"})

	_, err := stub.MockInvoke("t123", "RegistrarArrependimento", []string{pedidoID, "1503849607000"})
	if err == nil {
		t.Fatalf("Expected error ")
	}
}

func TestArrependimentoSuccess(t * testing.T) {
	fmt.Println("Entering TestArrependimentoErroDataAntes7Dias")
	attributes := make(map[string][]byte)
	stub := shim.NewCustomMockStub("mockStub", new(SaleContractChainCode), attributes)

	stub.MockInvoke("t123", "RegistrarPedido", []string{pedidoID, pedidoJson})

	stub.MockInvoke("t123", "RegistrarEntrega", []string{pedidoID, "1472313607000"})

	_, err := stub.MockInvoke("t123", "RegistrarArrependimento", []string{pedidoID, "1472313609000"})
	if err != nil {
		t.Fatalf("Not expected error ")
	}

	var pe Pedido
	ObterPedidoForTest(t, stub, pedidoID, &pe);
	if pe.Devolucao.MotivoDevolucao != 1 {
		t.Fatalf("Arrependimento not updated")
	}
}