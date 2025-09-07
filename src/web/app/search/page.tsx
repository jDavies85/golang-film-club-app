import { Container, Title } from '@mantine/core'
import MovieSearch from '@/components/MovieSearch'


export default function SearchPage() {
    return (
        <Container size="lg" py="xl">
            <Title order={2} mb="md">Search films</Title>
            <MovieSearch />
        </Container>
    )
}