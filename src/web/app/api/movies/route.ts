import { NextResponse } from 'next/server'

const API_BASE = process.env.API_BASE_URL ?? 'http://localhost:8080'

export const dynamic = 'force-dynamic'
export const revalidate = 0

export async function GET(req: Request) {
  const { searchParams } = new URL(req.url)
  const q = searchParams.get('q') ?? ''
  const page = searchParams.get('page') ?? '1'
  const pageSize = searchParams.get('pageSize') ?? '8'

  // Forward to Go API
  const url = new URL('/v1/movies/search', API_BASE)
  url.searchParams.set('query', q)
  url.searchParams.set('page', page)
  url.searchParams.set('pageSize', pageSize)

  const res = await fetch(url, { cache: 'no-store' })
  if (!res.ok) {
    return NextResponse.json({ error: `Upstream error ${res.status}` }, { status: 502 })
  }

  // If your Go API already returns {results,total,page,pageSize,totalPages}
  // you can just pass it through. If not, adapt/rename fields here.
  const data = await res.json()
  return NextResponse.json(data)
}
