'use client'


import { useCallback, useEffect, useMemo, useState } from 'react'
import { Button, Card, Grid, Group, Image, Pagination, Stack, Text, TextInput } from '@mantine/core'


export type Movie = {
    id: string
    title: string
    year: number
    posterUrl?: string
}


type ApiResponse = {
    query: string
    page: number
    pageSize: number
    total: number
    totalPages: number
    results: Movie[]
}


export default function MovieSearch() {
    const [query, setQuery] = useState('')
    const [page, setPage] = useState(1)
    const [pageSize] = useState(8)
    const [loading, setLoading] = useState(false)
    const [data, setData] = useState<ApiResponse | null>(null)
    const [error, setError] = useState<string | null>(null)
    const [submittedQuery, setSubmittedQuery] = useState<string>('')

    const canSearch = useMemo(() => query.trim().length > 0, [query])


    const runSearch = useCallback(async (q: string, p = 1) => {
        setLoading(true)
        setError(null)
        try {
            const res = await fetch(`/api/movies?q=${encodeURIComponent(q)}&page=${p}&pageSize=${pageSize}`, {
                cache: 'no-store',
            })
            if (!res.ok) throw new Error(`Request failed: ${res.status}`)
            const json: ApiResponse = await res.json()
            setData(json)
        } catch (e: any) {
            setError(e.message || 'Something went wrong')
        } finally {
            setLoading(false)
        }
    }, [pageSize])


    // whenever page changes, re-run with current query
    useEffect(() => {
        if (submittedQuery) {
            runSearch(submittedQuery, page)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [page, submittedQuery])


    return (
        <Stack gap="md">
            <form
                onSubmit={(e) => {
                    e.preventDefault()
                    setPage(1)
                    setSubmittedQuery(query.trim())
                    runSearch(query, 1)
                }}
            >
                <Group wrap="nowrap" gap="sm" align="end">
                    <TextInput
                        label="Search for a film"
                        placeholder="e.g. The Matrix"
                        value={query}
                        onChange={(e) => setQuery(e.currentTarget.value)}
                        w="100%"
                    />
                    <Button type="submit" disabled={!canSearch} loading={loading}>
                        Search
                    </Button>
                </Group>
            </form>


            {error && (
                <Card withBorder>
                    <Text c="red">{error}</Text>
                </Card>
            )}


            {data && (
                <Stack gap="sm">
                    <Text size="sm" c="dimmed">
                        Found {data.total} result{data.total === 1 ? '' : 's'}
                    </Text>


                    <Grid>
                        {data.results.map((m) => (
                            <Grid.Col key={m.id} span={{ base: 12, sm: 6, md: 3 }}>
                                <Card withBorder radius="lg" shadow="sm">
                                    <Card.Section>
                                        <Image
                                            src={m.posterUrl || 'https://placehold.co/600x900?text=Poster'}
                                            alt={`${m.title} poster`}
                                            h={220}
                                            fallbackSrc="https://placehold.co/600x900?text=No+Image"
                                            fit="cover"
                                        />
                                    </Card.Section>
                                    <Stack gap={4} mt="sm">
                                        <Text fw={600} lineClamp={2}>{m.title}</Text>
                                        <Text size="sm" c="dimmed">{m.year}</Text>
                                    </Stack>
                                </Card>
                            </Grid.Col>
                        ))}
                    </Grid>


                    {data.totalPages > 1 && (
                        <Group justify="center">
                            <Pagination value={page} onChange={setPage} total={data.totalPages} withEdges />
                        </Group>
                    )}
                </Stack>
            )}
        </Stack>
    )
}