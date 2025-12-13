import React, { useState, useEffect } from 'react';
import './EmployeeProfiling.css';
import { itemService } from './services/api';
import ItemForm from './components/ItemForm';
import ItemEditForm from './components/ItemEditForm';

function EmployeeProfiling() {
  const [employees, setEmployees] = useState([]);
  const [newEmployee, setNewEmployee] = useState({
    lastname: '',
    firstname: '',
    middlename: '',
    suffix: '',
    birthdate: '',
    sex: '',
    civil_status: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [editingId, setEditingId] = useState(null);
  const [editingEmployee, setEditingEmployee] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');

  const sexOptions = ['', 'Male', 'Female', 'Other'];
  const civilStatusOptions = ['', 'Single', 'Married', 'Divorced', 'Widowed'];

  useEffect(() => {
    fetchEmployees();
  }, []);

  const fetchEmployees = async () => {
    try {
      setLoading(true);
      const data = await itemService.getAllItems();
      setEmployees(data || []);
      setError('');
    } catch (err) {
      setError('Failed to fetch employees: ' + (err.response?.data?.error || err.message));
      console.error('Error fetching employees:', err);
    } finally {
      setLoading(false);
    }
  };

  const addEmployee = async (e) => {
    e.preventDefault();
    if (!newEmployee.lastname.trim() || !newEmployee.firstname.trim()) {
      setError('Lastname and Firstname are required');
      return;
    }
    try {
      setLoading(true);
      const data = await itemService.createItem(newEmployee);
      setEmployees([...employees, data]);
      setNewEmployee({
        lastname: '',
        firstname: '',
        middlename: '',
        suffix: '',
        birthdate: '',
        sex: '',
        civil_status: ''
      });
      setError('');
    } catch (err) {
      setError('Failed to add employee: ' + (err.response?.data?.error || err.message));
      console.error('Error adding employee:', err);
    } finally {
      setLoading(false);
    }
  };

  const deleteEmployee = async (id) => {
    const confirmDelete = window.confirm('Are you sure you want to delete this employee?');
    if (!confirmDelete) return;

    try {
      setLoading(true);
      await itemService.deleteItem(id);
      setEmployees(employees.filter(emp => emp.id !== id));
      setError('');
    } catch (err) {
      setError('Failed to delete employee: ' + err.message);
      console.error('Error deleting employee:', err);
    } finally {
      setLoading(false);
    }
  };

  const startEditing = (employee) => {
    setEditingId(employee.id);
    setEditingEmployee({ ...employee });
  };

  const cancelEditing = () => {
    setEditingId(null);
    setEditingEmployee(null);
  };

  const updateEmployee = async (id) => {
    if (!editingEmployee.lastname.trim() || !editingEmployee.firstname.trim()) {
      setError('Lastname and Firstname are required');
      return;
    }
    try {
      setLoading(true);
      const data = await itemService.updateItem(id, editingEmployee);
      setEmployees(employees.map(emp => emp.id === id ? data : emp));
      setEditingId(null);
      setEditingEmployee(null);
      setError('');
    } catch (err) {
      setError('Failed to update employee: ' + (err.response?.data?.error || err.message));
      console.error('Error updating employee:', err);
    } finally {
      setLoading(false);
    }
  };

  const filteredEmployees = employees.filter(emp =>
    emp.firstname?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    emp.lastname?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="profiling-container">
      <h2>👥 Employee Management</h2>
      <p className="subtitle">Create and manage employee profiles</p>
      
      {error && <div className="error">{error}</div>}
      
      <form onSubmit={addEmployee} className="add-form">
        <h3>Add New Employee</h3>
        <div className="form-grid">
          <input
            type="text"
            placeholder="Last Name *"
            value={newEmployee.lastname}
            onChange={(e) => setNewEmployee({...newEmployee, lastname: e.target.value})}
            required
          />
          <input
            type="text"
            placeholder="First Name *"
            value={newEmployee.firstname}
            onChange={(e) => setNewEmployee({...newEmployee, firstname: e.target.value})}
            required
          />
          <input
            type="text"
            placeholder="Middle Name"
            value={newEmployee.middlename}
            onChange={(e) => setNewEmployee({...newEmployee, middlename: e.target.value})}
          />
          <input
            type="text"
            placeholder="Suffix"
            value={newEmployee.suffix}
            onChange={(e) => setNewEmployee({...newEmployee, suffix: e.target.value})}
          />
          <input
            type="date"
            placeholder="Birthdate"
            value={newEmployee.birthdate}
            onChange={(e) => setNewEmployee({...newEmployee, birthdate: e.target.value})}
          />
          <select
            value={newEmployee.sex}
            onChange={(e) => setNewEmployee({...newEmployee, sex: e.target.value})}
          >
            {sexOptions.map((option, idx) => (
              <option key={idx} value={option}>
                {option || 'Select Sex'}
              </option>
            ))}
          </select>
          <select
            value={newEmployee.civil_status}
            onChange={(e) => setNewEmployee({...newEmployee, civil_status: e.target.value})}
          >
            {civilStatusOptions.map((option, idx) => (
              <option key={idx} value={option}>
                {option || 'Select Civil Status'}
              </option>
            ))}
          </select>
        </div>
        <button type="submit" disabled={loading} className="submit-btn">
          {loading ? 'Adding...' : 'Add Employee'}
        </button>
      </form>

      <div className="search-section">
        <input
          type="text"
          placeholder="🔍 Search by name..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="search-input"
        />
      </div>

      <div className="employees-section">
        <h3>Employee List ({filteredEmployees.length})</h3>
        {loading && <p>Loading...</p>}
        <div className="employees-table-container">
          <table className="employees-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Last Name</th>
                <th>First Name</th>
                <th>Middle Name</th>
                <th>Suffix</th>
                <th>Birthdate</th>
                <th>Sex</th>
                <th>Civil Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredEmployees.map(emp => (
                <tr key={emp.id}>
                  {editingId === emp.id ? (
                    <>
                      <td>{emp.id}</td>
                      <td><input value={editingEmployee.lastname} onChange={(e) => setEditingEmployee({...editingEmployee, lastname: e.target.value})} /></td>
                      <td><input value={editingEmployee.firstname} onChange={(e) => setEditingEmployee({...editingEmployee, firstname: e.target.value})} /></td>
                      <td><input value={editingEmployee.middlename} onChange={(e) => setEditingEmployee({...editingEmployee, middlename: e.target.value})} /></td>
                      <td><input value={editingEmployee.suffix} onChange={(e) => setEditingEmployee({...editingEmployee, suffix: e.target.value})} /></td>
                      <td><input type="date" value={editingEmployee.birthdate} onChange={(e) => setEditingEmployee({...editingEmployee, birthdate: e.target.value})} /></td>
                      <td>
                        <select value={editingEmployee.sex} onChange={(e) => setEditingEmployee({...editingEmployee, sex: e.target.value})}>
                          {sexOptions.map((option, idx) => <option key={idx} value={option}>{option || 'Select'}</option>)}
                        </select>
                      </td>
                      <td>
                        <select value={editingEmployee.civil_status} onChange={(e) => setEditingEmployee({...editingEmployee, civil_status: e.target.value})}>
                          {civilStatusOptions.map((option, idx) => <option key={idx} value={option}>{option || 'Select'}</option>)}
                        </select>
                      </td>
                      <td>
                        <button onClick={() => updateEmployee(emp.id)} className="save-btn">Save</button>
                        <button onClick={cancelEditing} className="cancel-btn">Cancel</button>
                      </td>
                    </>
                  ) : (
                    <>
                      <td>{emp.id}</td>
                      <td>{emp.lastname}</td>
                      <td>{emp.firstname}</td>
                      <td>{emp.middlename}</td>
                      <td>{emp.suffix}</td>
                      <td>{emp.birthdate}</td>
                      <td>{emp.sex}</td>
                      <td>{emp.civil_status}</td>
                      <td>
                        <button onClick={() => startEditing(emp)} className="edit-btn">Edit</button>
                        <button onClick={() => deleteEmployee(emp.id)} className="delete-btn">Delete</button>
                      </td>
                    </>
                  )}
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

export default EmployeeProfiling;
