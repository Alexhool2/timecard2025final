import { useState } from "react"
import { Badge, Button, Form, ListGroup, Spinner } from "react-bootstrap"

const API_URL = import.meta.env.VITE_API_URL
interface User {
  id: number
  firstName: string
  lastName: string
  userName: string
  email: string
}

const UserSearch = ({
  onSelectUser,
}: {
  onSelectUser: (id: number) => void
}) => {
  const [searchTerm, setSearchTerm] = useState("")
  const [users, setUsers] = useState<User[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [selectedUsers, setSelectedUsers] = useState<User[]>([])

  const handleSearch = async (term: string) => {
    setSearchTerm(term)
    if (term.length < 2) {
      setUsers([])
      return
    }

    setIsLoading(true)
    try {
      const response = await fetch(`${API_URL}/api/v1/users?search=${term}`, {
        method: "GET",
        credentials: "include",
      })
      const data = await response.json()
      setUsers(data || [])
    } catch (error) {
      console.log("Error on searching users", error)
      setUsers([])
    } finally {
      setIsLoading(false)
    }
  }

  //add selected user
  const handleSelectedUser = (user: User) => {
    if (selectedUsers.length > 1) {
      return console.log("only possible to select 1 user by search")
    }
    if (!selectedUsers.some((u) => u.id === user.id)) {
      setSelectedUsers([...selectedUsers, user])
      onSelectUser(user.id)
    }
    setUsers([])
    setSearchTerm("")
  }
  //VERIFICAR PROBLEMAS NA BUSCA DOS USUARIOS E HORARIOS RETORNADOS
  const handleRemoveUser = (id: number) => {
    setSelectedUsers(selectedUsers.filter((user) => user.id !== id))
    setUsers([])
    setSearchTerm("")
    localStorage.clear()
    onSelectUser(0)
  }

  return (
    <div>
      <Form.Control
        type="text"
        placeholder="Insert name, userName, or email"
        value={searchTerm}
        onChange={(e) => handleSearch(e.target.value)}
        disabled={selectedUsers.length > 0}
      />

      {/* Exibe o spinner abaixo do campo */}
      {isLoading && <Spinner animation="border" className="mt-2" />}

      {/* Lista de usuários encontrados */}
      {users.length > 0 && (
        <ListGroup className="mt-2">
          {users.map((user) => (
            <ListGroup.Item
              key={user.id}
              action
              onClick={() => handleSelectedUser(user)}
            >
              {user.userName} ({user.email}) ({user.firstName} {user.lastName})
            </ListGroup.Item>
          ))}
        </ListGroup>
      )}
      {users.length === 0 && searchTerm.length > 2 && (
        <p className="text-danger mt-2">No users found</p>
      )}
      {/* selected users return */}
      {selectedUsers.length > 0 && (
        <div className="mt-3">
          <h5>Selected user:</h5>
          {selectedUsers.map((user) => (
            <Badge
              key={user.id}
              bg="dark"
              className="me-2 p-2 d-inline-flex align-items-center mt-1"
            >
              {user.userName} ({user.email})
              <Button
                variant="light"
                size="sm"
                className="ms-2"
                onClick={() => handleRemoveUser(user.id)}
              >
                X
              </Button>
            </Badge>
          ))}
        </div>
      )}
    </div>
  )
}

export default UserSearch
