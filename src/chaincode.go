package main

import (
	"errors"
	"fmt"
	// "encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("mylogger")

type Troca struct {
	// 1 Arrependimento
	// 2 Defeituoso
	MotivoTroca            int32         `json:"motivoTroca"`
	// 1 devolucao pagamento
	// 2 abatimento
	// 3 por produto
	OpcaoTroca			   int32		 `json:"opcaoTroca"`
	Data              	   int64         `json:"data"`
}


type Devolucao struct {
	// 1 Arrependimento
	// 2 Defeituoso
	// 3 Caso Fortuito
	MotivoDevolucao        int32         `json:"devolvido"`
	ComplementoMotivoDevolucao   string  `json:"complementoMotivoDevolucao"`
	Data              	   int64         `json:"data"`
}

type Pedido struct {
	ID                     string        `json:"id"`
	CPFCliente			   string 		 `json:"cpf"`				
	DescricaoItens         string        `json:"descricaoItens"`
	ItensId                string        `json:"itensId"`
	DataVenda              int64         `json:"dataVenda"`
	DataEntrega            int64         `json:"dataEntrega"`
	Devolucao              Devolucao     `json:"devolucao"`
	Troca              	   Troca         `json:"troca"`
}

//CONTRACT
type SaleContractChainCode struct {
}

func main() {
	err := shim.Start(new(SaleContractChainCode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

func (t *SaleContractChainCode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("init")
    return nil, nil
}
 
func (t *SaleContractChainCode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if function == "ObterPedido" {
		return ObterPedido(stub, args)
	} else {
		return nil, errors.New(" Unknow query method ")
	} 
	return nil, nil
}
 
func (t *SaleContractChainCode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
    if function == "RegistrarPedido" {
		return RegistrarPedido(stub, args)
	} else {
		return nil, errors.New(" Unknow invoke method ")
	} 
	return nil, nil
}

func RegistrarPedido(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	logger.Debug("Entering RegistrarPedido")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for loan application creation")
	}

	var pedidoID = args[0]
	var pedidoInput = args[1]

	//TODO validar schema do json
	//TODO validar permissao para criacao de pedido
	
	err := stub.PutState(pedidoID, []byte(pedidoInput))
	if err != nil {
		logger.Error("Could not save pedido to ledger", err)
		return nil, err
	}
	logger.Info("Successfully saved Pedido");

	return []byte(pedidoInput), nil
}

func ObterPedido(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	logger.Debug("Entering ObterPedido")

	if len(args) < 1 {
		logger.Error("Invalid number of arguments")
		return nil, errors.New("Missing pedido ID")
	} 

	var pedidoId = args[0]
	bytes, err := stub.GetState(pedidoId)
	if err != nil {
		logger.Error("Could not fetch loan application with id "+pedidoId+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}
