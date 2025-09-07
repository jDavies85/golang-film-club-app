import { NextResponse } from 'next/server'
import { MOVIES } from '@/lib/mock/movies'


export const dynamic = 'force-dynamic';
export const revalidate = 0;

export async function GET(req: Request) {
    const { searchParams } = new URL(req.url)
    const q = (searchParams.get('q') || '').toLowerCase().trim()
    const page = Math.max(1, Number(searchParams.get('page') || 1))
    const pageSize = Math.min(50, Math.max(1, Number(searchParams.get('pageSize') || 8)))


    // filter by title match; expand later with TMDB fields
    const filtered = q
        ? MOVIES.filter(m => m.title.toLowerCase().includes(q))
        : MOVIES


    const total = filtered.length
    const totalPages = Math.max(1, Math.ceil(total / pageSize))
    const clampedPage = Math.min(page, totalPages)
    const startIdx = (clampedPage - 1) * pageSize
    const results = filtered.slice(startIdx, startIdx + pageSize)


    return NextResponse.json({
        query: q,
        page: clampedPage,
        pageSize,
        total,
        totalPages,
        results,
    })
}