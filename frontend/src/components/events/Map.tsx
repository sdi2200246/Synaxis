import { MapContainer, TileLayer, Marker, Popup } from 'react-leaflet'
import L from 'leaflet'
import 'leaflet/dist/leaflet.css'

import markerIcon from 'leaflet/dist/images/marker-icon.png'
import markerShadow from 'leaflet/dist/images/marker-shadow.png'
import markerRetina from 'leaflet/dist/images/marker-icon-2x.png'

delete (L.Icon.Default.prototype as any)._getIconUrl

L.Icon.Default.mergeOptions({
  iconUrl: markerIcon,
  shadowUrl: markerShadow,
  iconRetinaUrl: markerRetina,
})

interface EventMapProps {
  lat: number
  lng: number
  venueName: string
}

export function EventMap({ lat, lng, venueName }: EventMapProps) {
  return (
    <MapContainer
      center={[lat, lng]}
      zoom={15}
      className="event-map"
      scrollWheelZoom={false}
    >
      <TileLayer
        url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
        attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
      />
      <Marker position={[lat, lng]}>
        <Popup>{venueName}</Popup>
      </Marker>
    </MapContainer>
  )
}