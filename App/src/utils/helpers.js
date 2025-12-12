export const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

export const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

export const getStatusBadge = (status) => {
  const statusMap = {
    draft: { class: 'badge-info', text: 'Draft' },
    pending_verification: { class: 'badge-warning', text: 'Pending' },
    verified: { class: 'badge-success', text: 'Verified' },
    rejected: { class: 'badge-danger', text: 'Rejected' },
  }
  return statusMap[status] || { class: 'badge-info', text: status }
}

export const truncate = (str, length = 50) => {
  if (!str) return ''
  return str.length > length ? str.substring(0, length) + '...' : str
}

export const getFileUrl = (filename) => {
  return `${import.meta.env.VITE_API_BASE_URL.replace('/api/v1', '')}/uploads/${filename}`
}
