package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("mylogger")

type Troca struct {
	// 1 Arrependimento
	// 2 Defeituoso
	MotivoTroca            int         `json:"motivoTroca"`
	// 1 devolucao pagamento
	// 2 abatimento
	// 3 por produto
	OpcaoTroca			   int		 `json:"opcaoTroca"`
	Data              	   int64         `json:"data"`
}


type Devolucao struct {
	// 1 Arrependimento
	// 2 Defeituoso
	// 3 Caso Fortuito
	MotivoDevolucao        int         `json:"devolvido"`
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
	} 
	if function == "RegistrarEntrega" {
		return RegistrarEntrega(stub, args)
	} 
	if function == "Arrependimento" {
		return Arrependimento(stub, args)
	} else {
		return nil, errors.New(" Unknow invoke method ")
	} 
	return nil, nil
}

func AtualizarPedido( stub shim.ChaincodeStubInterface, id string, fn func(p *Pedido) error ) ([]byte, error){
	
	bytes, err := ObterPedido(stub, []string{id});
	if err != nil {
		logger.Error("Invalid format " + id, err)
		return nil, err
	}

	var pe Pedido
	err = json.Unmarshal(bytes, &pe)
	if err != nil {
		logger.Error("Invalid format atualizar unmarshal " + string(bytes[:]), err)
		return nil, errors.New(" Invalid json format ")
	}

	err = fn(&pe)
	if err != nil {
		logger.Error("validation error in update function ", err)
		return nil, err
	}	

	bytes, err = json.Marshal(&pe)
	if err != nil {
		logger.Error("Could not marshal Pedido: update", err)
		return nil, err
	}
	
	err = stub.PutState(id, bytes)
	if err != nil {
		logger.Error("Could not update pedido to ledger", err)
		return nil, err
	}
	logger.Info("Successfully updated Pedido");

	return bytes, nil

}

func RegistrarEntrega (stub shim.ChaincodeStubInterface, args []string ) ([]byte, error) {
	
	logger.Debug("Entering RegistrarEntrega")
	
	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for Registrar Entrega")
	}
	var pedidoID = args[0]
	var dataEntrega = args[1]
	dataEntregaLong, err := strconv.ParseInt(dataEntrega, 10, 64);
	if err != nil {
		logger.Error("Invalid timestamp value")
		return nil, errors.New("Invalid timestamp value")	
	}

	fn := func(p *Pedido) error {		
		p.DataEntrega = dataEntregaLong
		return nil
	}

	return AtualizarPedido(stub, pedidoID, fn)
}

func Arrependimento( stub shim.ChaincodeStubInterface, args []string )  ([]byte, error) {
	
	logger.Debug("Entering Arrependimento")
	
	if len(args) < 3 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast tree arguments for Arrependimento")
	}

	var pedidoID = args[0]
	var codigoArrependimento = args[1]
	var dataDevolucao = args[2]
	codigoArrependimentoInt, err := strconv.Atoi(codigoArrependimento);
	if err != nil {
		logger.Error("Invalid number value")
		return nil, errors.New("Invalid number value")	
	}
	dataDevolucaoLong, err := strconv.ParseInt(dataDevolucao, 10, 64);
	if err != nil {
		logger.Error("Invalid timestamp value")
		return nil, errors.New("Invalid timestamp value")	
	}

	fn := func(p *Pedido) error {		
		if p.DataEntrega == 0 {
			return errors.New("Product that was not delivered can not be returned")
		}
		//se for maior que 7 dias a diferenca nao deixa se arrepender
		if( dataDevolucaoLong - p.DataEntrega > 604800000  ) {
			logger.Error(codigoArrependimentoInt)
			return errors.New("Time of regret exceeded")
		}
		p.Devolucao.MotivoDevolucao = codigoArrependimentoInt;
		p.Devolucao.Data = dataDevolucaoLong;
		return nil
	}
	return AtualizarPedido(stub, pedidoID, fn)
}

func RegistrarPedido(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	logger.Debug("Entering RegistrarPedido")

	if len(args) < 2 {
		logger.Error("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for loan application creation")
	}

	var pedidoID = args[0]
	var pedidoInput = args[1]

	var pe Pedido
	err := json.Unmarshal([]byte(pedidoInput), &pe)
	if err != nil {
		logger.Error("Invalid format", err)
		return nil, errors.New(" Invalid json format ")
	}

	pe.ID = pedidoID

	peBytes, err := json.Marshal(&pe)
	if err != nil {
		logger.Error("Could not marshal Pedido", err)
		return nil, err
	}

	//TODO validar schema do json
	//TODO validar permissao para criacao de pedido
	
	err = stub.PutState(pedidoID, peBytes)
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
