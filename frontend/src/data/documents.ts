import type { DocumentItem } from './types'

/* Документы: лицензии, устав, сертификаты.
   url — плейсхолдер '#', подставь реальные ссылки на файлы. */
export const documents: DocumentItem[] = [
  {
    id: '1',
    title: 'Лицензия на образовательную деятельность',
    description: 'Действующая лицензия с приложениями.',
    fileType: 'PDF',
    fileSize: '1,4 МБ',
    url: '#',
  },
  {
    id: '2',
    title: 'Устав центра',
    description: 'Учредительный документ организации.',
    fileType: 'PDF',
    fileSize: '820 КБ',
    url: '#',
  },
  {
    id: '3',
    title: 'Свидетельство о регистрации',
    description: 'Документ о государственной регистрации.',
    fileType: 'PDF',
    fileSize: '540 КБ',
    url: '#',
  },
  {
    id: '4',
    title: 'Правила внутреннего распорядка',
    description: 'Порядок посещения занятий и правила центра.',
    fileType: 'PDF',
    fileSize: '310 КБ',
    url: '#',
  },
  {
    id: '5',
    title: 'Договор об оказании услуг (образец)',
    description: 'Типовой договор для ознакомления перед записью.',
    fileType: 'PDF',
    fileSize: '260 КБ',
    url: '#',
  },
  {
    id: '6',
    title: 'Политика обработки персональных данных',
    description: 'Как мы храним и защищаем ваши данные.',
    fileType: 'PDF',
    fileSize: '190 КБ',
    url: '#',
  },
]
