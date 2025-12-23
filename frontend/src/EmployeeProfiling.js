import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './EmployeeProfiling.css';
import { itemService } from './services/api';
import ItemForm from './components/ItemForm';
import ItemEditForm from './components/ItemEditForm';

function EmployeeProfiling() {
  const navigate = useNavigate();
  const [employees, setEmployees] = useState([]);
  const [newEmployee, setNewEmployee] = useState({
    caseNo: '',
    hospitalNo: '',
    lastname: '',
    firstname: '',
    middlename: '',
    suffix: '',
    birthdate: '',
    age: '',
    room: '',
    admissionDate: '',
    dischargeDate: '',
    sex: '',
    height: '',
    weight: '',
    complaint: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [searchTerm, setSearchTerm] = useState('');
  const [editItem, setEditItem] = useState(null);
  const [editLoading, setEditLoading] = useState(false);

  const sexOptions = ['', 'Male', 'Female', 'Other'];
  const civilStatusOptions = ['', 'Single', 'Married', 'Divorced', 'Widowed'];

  // Input sanitizers
  const onlyDigits = (val) => (val || '').toString().replace(/\D+/g, '');
  const onlyLetters = (val) => (val || '').toString().replace(/[^a-zA-Z\s]+/g, '');

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
        caseNo: '',
        hospitalNo: '',
        lastname: '',
        firstname: '',
        middlename: '',
        suffix: '',
        birthdate: '',
        age: '',
        room: '',
        admissionDate: '',
        dischargeDate: '',
        sex: '',
        height: '',
        weight: '',
        complaint: ''
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

  const openEdit = (emp) => {
    setEditItem({...emp});
  };

  const handleEditChange = (updated) => {
    setEditItem(updated);
  };

  const saveEdit = async () => {
    if (!editItem) return;
    try {
      setEditLoading(true);
      const updated = await itemService.updateItem(editItem.id, editItem);
      setEmployees(employees.map(e => e.id === updated.id ? updated : e));
      setEditItem(null);
      setError('');
    } catch (err) {
      setError('Failed to update employee: ' + (err.response?.data?.error || err.message));
      console.error('Error updating employee:', err);
    } finally {
      setEditLoading(false);
    }
  };

  const cancelEdit = () => {
    setEditItem(null);
  };

  const viewPatient = (patient) => {
    navigate(`/patient/${patient.id}`);
  };

  const filteredEmployees = employees.filter(emp =>
    emp.firstname?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    emp.lastname?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="profiling-container">
      <h2>👥 Patient Management</h2>
      <p className="subtitle">Create and manage patient profiles</p>
      
      {error && <div className="error">{error}</div>}
      
      <form onSubmit={addEmployee} className="add-form">
        <h3>Add New Patient</h3>
        <div className="form-grid">
          <div className="form-field">
            <label>Case No. *</label>
            <input
              name="caseNo"
              type="text"
              value={newEmployee.caseNo}
              onChange={(e) => setNewEmployee({...newEmployee, caseNo: onlyDigits(e.target.value)})}
              required
            />
          </div>
          <div className="form-field">
            <label>Hospital No.</label>
            <input
              name="hospitalNo"
              type="text"
              value={newEmployee.hospitalNo}
              onChange={(e) => setNewEmployee({...newEmployee, hospitalNo: onlyDigits(e.target.value)})}
            />
          </div>
          <div className="form-field">
            <label>Last Name *</label>
            <input
              name="lastname"
              type="text"
              value={newEmployee.lastname}
              onChange={(e) => setNewEmployee({...newEmployee, lastname: onlyLetters(e.target.value)})}
              required
            />
          </div>
          <div className="form-field">
            <label>First Name *</label>
            <input
              name="firstname"
              type="text"
              value={newEmployee.firstname}
              onChange={(e) => setNewEmployee({...newEmployee, firstname: onlyLetters(e.target.value)})}
              required
            />
          </div>
          
          <div className="form-field">
            <label>Middle Name</label>
            <input
              name="middlename"
              type="text"
              value={newEmployee.middlename}
              onChange={(e) => setNewEmployee({...newEmployee, middlename: onlyLetters(e.target.value)})}
            />
          </div>
          <div className="form-field">
            <label>Suffix</label>
            <input
              name="suffix"
              type="text"
              value={newEmployee.suffix}
              onChange={(e) => setNewEmployee({...newEmployee, suffix: onlyLetters(e.target.value)})}
            />
          </div>
          <div className="form-field">
            <label>Birthdate *</label>
            <input
              name="birthdate"
              type="date"
              value={newEmployee.birthdate}
              onChange={(e) => setNewEmployee({...newEmployee, birthdate: e.target.value})}
              required
            />
          </div>
          <div className="form-field">
            <label>Age</label>
            <input
              name="age"
              type="text"
              value={newEmployee.age}
              onChange={(e) => setNewEmployee({...newEmployee, age: onlyDigits(e.target.value)})}
            />
          </div>
          
          <div className="form-field">
            <label>Room</label>
            <input
              name="room"
              type="text"
              value={newEmployee.room}
              onChange={(e) => setNewEmployee({...newEmployee, room: onlyDigits(e.target.value)})}
            />
          </div>
          <div className="form-field">
            <label>Admission Date</label>
            <input
              name="admissionDate"
              type="datetime-local"
              value={newEmployee.admissionDate}
              onChange={(e) => setNewEmployee({...newEmployee, admissionDate: e.target.value})}
            />
          </div>
          <div className="form-field">
            <label>Discharge Date</label>
            <input
              name="dischargeDate"
              type="datetime-local"
              value={newEmployee.dischargeDate}
              onChange={(e) => setNewEmployee({...newEmployee, dischargeDate: e.target.value})}
            />
          </div>
          <div className="form-field">
            <label>Sex</label>
            <select
              name="sex"
              value={newEmployee.sex}
              onChange={(e) => setNewEmployee({...newEmployee, sex: e.target.value})}
            >
              <option value="">Select Sex</option>
              <option value="Male">Male</option>
              <option value="Female">Female</option>
            </select>
          </div>
          
          <div className="form-field">
            <label>Height (cm)</label>
            <input
              name="height"
              type="text"
              value={newEmployee.height}
              onChange={(e) => setNewEmployee({...newEmployee, height: onlyDigits(e.target.value)})}
            />
          </div>
          <div className="form-field">
            <label>Weight (kg)</label>
            <input
              name="weight"
              type="text"
              value={newEmployee.weight}
              onChange={(e) => setNewEmployee({...newEmployee, weight: onlyDigits(e.target.value)})}
            />
          </div>
          <div className="form-field">
            <label>Complaint</label>
            <input
              name="complaint"
              type="text"
              value={newEmployee.complaint}
              onChange={(e) => setNewEmployee({...newEmployee, complaint: onlyLetters(e.target.value)})}
            />
          </div>
          <div className="form-field">
            <label>&nbsp;</label>
            <button type="submit" disabled={loading} className="submit-btn">
              {loading ? 'Adding...' : 'Add Patient'}
            </button>
          </div>
        </div>
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
      {/* Edit Modal */}
      {editItem && (
        <div className="modal-overlay">
          <div className="modal-content">
            <div className="modal-header">
              <h3>Edit Patient</h3>
              <button className="close-btn" onClick={cancelEdit}>&times;</button>
            </div>
            <div className="patient-details">
              <ItemEditForm
                item={editItem}
                onChange={handleEditChange}
                onSave={saveEdit}
                onCancel={cancelEdit}
                loading={editLoading}
                sexOptions={sexOptions}
                civilStatusOptions={civilStatusOptions}
              />
            </div>
          </div>
        </div>
      )}

      <div className="employees-section">
        <h3>Patient List ({filteredEmployees.length})</h3>
        {loading && <p>Loading...</p>}
        <div className="employees-table-container">
          <table className="employees-table">
            <thead>
              <tr>
                <th>Case No.</th>
                <th>Full Name</th>
                <th>Birthdate</th>
                <th>Room</th>
                <th>Admission Date</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredEmployees.map(emp => (
                <tr key={emp.id}>
                  <td>{emp.caseNo}</td>
                  <td>{`${emp.lastname}, ${emp.firstname} ${emp.middlename || ''} ${emp.suffix || ''}`.trim()}</td>
                  <td>{emp.birthdate}</td>
                  <td>{emp.room}</td>
                  <td>{emp.admissionDate ? new Date(emp.admissionDate).toLocaleString() : ''}</td>
                  <td>
                    <button onClick={() => openEdit(emp)} className="edit-btn">Edit</button>
                    <button onClick={() => viewPatient(emp)} className="view-btn">View</button>
                    <button onClick={() => deleteEmployee(emp.id)} className="delete-btn">Delete</button>
                  </td>
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
