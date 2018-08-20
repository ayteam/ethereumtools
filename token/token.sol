pragma solidity ^0.4.11;

contract Token {

    string public name = "";      //  token name
    string public symbol = "";           //  token symbol
    uint256 public decimals = 2;            //  token digit
    uint256 public totalSupply = 0;

    mapping (address => uint256) public balanceOf;
    mapping (address => bool) public lockedAccount;
    mapping (address => mapping (address => uint256)) public allowance;

    bool public stopped = false;
    address owner = 0x0;

    modifier isOwner {
        assert(owner == msg.sender);
        _;
    }

    modifier isRunning {
        assert (!stopped);
        _;
    }

    modifier validAddress {
        assert(0x0 != msg.sender);
        _;
    }

    modifier unLocked {
        assert(lockedAccount[msg.sender] != true);
        _;
    }

    constructor (string _name, string _symbol, uint256 _decimals, uint256 _total, address _addressFounder) public {
        owner = msg.sender;
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        totalSupply = _total;

        balanceOf[_addressFounder] = totalSupply;
        emit Transfer(0x0, _addressFounder, totalSupply);
    }

    function increate(address _addressFounder, uint256 _value) isRunning unLocked isOwner public {
        require(balanceOf[_addressFounder] + _value >= balanceOf[_addressFounder]);
        require(totalSupply + _value > totalSupply);
        balanceOf[_addressFounder] += _value;
        totalSupply += _value;
        emit Transfer(0x0, _addressFounder, _value);
    }

    function transfer(address _to, uint256 _value) isRunning validAddress unLocked public returns (bool success) {
        require(balanceOf[msg.sender] >= _value);
        require(balanceOf[_to] + _value >= balanceOf[_to]);
        balanceOf[msg.sender] -= _value;
        balanceOf[_to] += _value;
        emit Transfer(msg.sender, _to, _value);
        return true;
    }

    function transferFrom(address _from, address _to, uint256 _value) isRunning validAddress unLocked public returns (bool success) {
        require(balanceOf[_from] >= _value);
        require(balanceOf[_to] + _value >= balanceOf[_to]);
        require(allowance[_from][msg.sender] >= _value);
        balanceOf[_to] += _value;
        balanceOf[_from] -= _value;
        allowance[_from][msg.sender] -= _value;
        emit Transfer(_from, _to, _value);
        return true;
    }

    function approve(address _spender, uint256 _value) isRunning validAddress unLocked public returns (bool success) {
        require(_value > 0);
        allowance[msg.sender][_spender] = _value;
        emit Approval(msg.sender, _spender, _value);
        return true;
    }

    function stop() isOwner public {
        stopped = true;
    }

    function start() isOwner public {
        stopped = false;
    }

    function setName(string _name) isOwner public {
        name = _name;
    }

    function lockAccount(address _account) isOwner public {
        lockedAccount[_account] = true;
    }

    function unlockAccount(address _account) isOwner public {
        delete lockedAccount[_account];
    }

    function isLocked(address _account) public returns (bool success) {
        if (lockedAccount[_account] == true)
            emit IsLocked(_account, true);
        else
            emit IsLocked(_account, false);
        return true;
    }

    function burn(uint256 _value) public {
        require(balanceOf[msg.sender] >= _value);
        balanceOf[msg.sender] -= _value;
        balanceOf[0x0] += _value;
        emit Transfer(msg.sender, 0x0, _value);
    }

    event Transfer(address indexed _from, address indexed _to, uint256 _value);
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);
    event IsLocked(address indexed _account, bool lock);
}
