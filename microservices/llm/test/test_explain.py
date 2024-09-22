import unittest
from unittest.mock import patch, MagicMock
from handlers import explain


class TestExplain(unittest.TestCase):

    @patch('handlers.explain.client.chat.completions.create')
    def test_synthesize_major_topics(self, mock_create):
        mock_response = MagicMock()
        mock_response.choices[0].message.content = "Topic 1\nTopic 2"
        mock_create.return_value = mock_response

        result = explain.synthesize_major_topics("Sample text")
        self.assertEqual(result, ["Topic 1", "Topic 2"])
        mock_create.assert_called_once()

    @patch('handlers.explain.client.chat.completions.create')
    def test_synthesize_major_concepts(self, mock_create):
        mock_response = MagicMock()
        mock_response.choices[0].message.content = "Concept 1\nConcept 2"
        mock_create.return_value = mock_response

        result = explain.synthesize_major_concepts("Sample topic")
        self.assertEqual(result, ["Concept 1", "Concept 2"])
        mock_create.assert_called_once()

    @patch('handlers.explain.client.chat.completions.create')
    def test_generate_layered_explanation(self, mock_create):
        mock_response = MagicMock()
        mock_response.choices[0].message.content = "Layer explanation"
        mock_create.return_value = mock_response

        result = explain.generate_layered_explanation("Sample concept")
        self.assertEqual(len(result), 3)
        self.assertEqual(result[0], "Layer explanation")
        self.assertEqual(mock_create.call_count, 3)

    @patch('handlers.explain.synthesize_major_topics')
    @patch('handlers.explain.synthesize_major_concepts')
    @patch('handlers.explain.generate_layered_explanation')
    def test_process_text(self, mock_generate, mock_concepts, mock_topics):
        mock_topics.return_value = ["Topic 1"]
        mock_concepts.return_value = ["Concept 1"]
        mock_generate.return_value = ["Layer 1", "Layer 2", "Layer 3"]

        result = explain.process_text("Sample text")
        expected_result = [
            {
                "topic": "Topic 1",
                "concepts": [
                    {
                        "concept": "Concept 1",
                        "layers": [
                            {"layer_1": "Layer 1"},
                            {"layer_2": "Layer 2"},
                            {"layer_3": "Layer 3"}
                        ]
                    }
                ]
            }
        ]
        self.assertEqual(result, expected_result)


if __name__ == '__main__':
    unittest.main()